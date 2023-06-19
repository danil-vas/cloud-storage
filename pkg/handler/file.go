package handler

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

type jsonResponse struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Status string `json:"status"`
}

var Access = make(map[int]int)

func deleteTempFile(str string) {
	for true {
		os.Remove(str)
		_, err := os.Stat(str)
		if err == nil {
			time.Sleep(time.Second * 1)
		} else {
			return
		}
	}

}

// @Summary Share Object
// @Security ApiKeyAuth
// @Tags User
// @Description share object
// @ID share-object
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /api/user/id [post]
func (h *Handler) shareObject(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	idObject, _ := c.GetPostForm("idObject")
	objectId, err := strconv.Atoi(idObject)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	flag, err := h.services.CheckAccessToObject(userId, objectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if flag == false {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}
	shareUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	fmt.Println("WORK")
	Access[objectId] = shareUserId
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "complete",
	})
}

// @Summary Download File
// @Security ApiKeyAuth
// @Tags File
// @Description download file
// @ID download-file
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /api/file/{id} [get]
func (h *Handler) downloadFile(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	objectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	flag, err := h.services.CheckAccessToObject(userId, objectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if flag == false {
		i, ok := Access[objectId]
		if !(ok && i == userId) {
			c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
			return
		}
	}
	typeFile, err := h.services.GetTypeObject(objectId)
	path, err := h.services.File.PathUploadFile(userId, objectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	path = "temp/" + path
	if typeFile == "file" {
		file, err := os.Open(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		originalFileName, err := h.services.File.OriginalFileName(objectId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		fmt.Println(path)
		contentDisposition := fmt.Sprintf("attachment; filename=%s", originalFileName)
		c.Header("Content-Disposition", contentDisposition)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.File(path)
	} else {
		zipFilePath := path + "Archive.zip"
		zipFile, err := os.Create(zipFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		defer zipFile.Close()

		archive := zip.NewWriter(zipFile)
		rootPath := "./" + path
		err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			relPath, err := filepath.Rel(rootPath, path)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			name := ""
			for i := len(relPath) - 1; i != -1; i-- {
				if string(relPath[i]) != `\` {
					name = string(relPath[i]) + name
				} else {
					break
				}
			}
			nameFile, err := h.services.File.OriginalFileNameThroughServerName(name)
			if err != nil {
				return err
			}
			zipFile, err := archive.Create(nameFile)
			if err != nil {
				return err
			}
			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		archive.Close()
		file, err := os.Open(zipFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		originalFileName, err := h.services.File.OriginalFileName(objectId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		originalFileName += "Archive.zip"
		contentDisposition := fmt.Sprintf("attachment; filename=%s", originalFileName)
		c.Header("Content-Disposition", contentDisposition)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.File(zipFilePath)
		go deleteTempFile(zipFilePath)
	}
}

// @Summary Delete File
// @Security ApiKeyAuth
// @Tags File
// @Description delete file
// @ID delete-file
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /api/file/{id} [delete]
func (h *Handler) deleteFile(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	objectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	flag, err := h.services.CheckAccessToObject(userId, objectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if flag == false {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}
	path, err := h.services.File.PathUploadFile(userId, objectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	path = "temp/" + path
	err = os.Remove(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	err = h.services.File.DeleteFile(userId, objectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "complete",
	})
}

// @Summary Upload File
// @Security ApiKeyAuth
// @Tags File
// @Description upload file
// @ID upload-file
// @Accept  mpfd
// @Produce  json
// @Success 200 {object} jsonResponse
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /api/file/{id} [post]
func (h *Handler) uploadFile(ctx *gin.Context) {
	userId, err := getUsersId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error authorization"})
		return
	}
	objectId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	flag, err := h.services.CheckAccessToObject(userId, objectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if flag == false {
		ctx.JSON(http.StatusForbidden, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	typeObject, err := h.services.File.GetTypeObject(objectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if typeObject == "file" {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "unable to upload file to file",
		})
		return
	}
	path, err := h.services.File.PathUploadFile(userId, objectId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	form, _ := ctx.MultipartForm()
	files := form.File["myFile"]
	str := make([]string, 0)
	resp := make([]jsonResponse, 0)
	for _, file := range files {
		var infoFile jsonResponse
		size := file.Size
		availableMemory, err := h.services.File.GetAvailableMemory(userId)
		if err != nil {
			break
		}
		if int64(availableMemory) < size {
			infoFile.Id = 0
			infoFile.Name = file.Filename
			infoFile.Size = file.Size
			infoFile.Status = "there is no available space for the user to upload the file"
			resp = append(resp, infoFile)
			break
		}
		str = append(str, file.Filename)
		nameFile := file.Filename
		extension := ""
		for i := len(nameFile) - 1; i != 0; i-- {
			if string(nameFile[i]) != "." {
				extension = string(nameFile[i]) + extension
			} else {
				break
			}
		}
		tempFile, err := os.CreateTemp("temp/"+path, "*."+extension)
		if err != nil {
			infoFile.Id = 0
			infoFile.Name = file.Filename
			infoFile.Size = file.Size
			infoFile.Status = err.Error()
			resp = append(resp, infoFile)
			break
		}
		serverNameFile := filepath.Base(tempFile.Name())
		info, _ := tempFile.Stat()
		timeFile := info.Sys().(*syscall.Win32FileAttributeData).LastAccessTime
		fileAccessTime := time.Unix(0, timeFile.Nanoseconds())
		defer tempFile.Close()
		readerFile, _ := file.Open()
		_, err = io.Copy(tempFile, readerFile)
		if err != nil {
			infoFile.Id = 0
			infoFile.Name = file.Filename
			infoFile.Size = file.Size
			infoFile.Status = err.Error()
			resp = append(resp, infoFile)
			break
		}
		idFile, err := h.services.AddUploadFileToUser(userId, objectId, file.Filename, serverNameFile, file.Size, fileAccessTime)
		if err != nil {
			infoFile.Id = 0
			infoFile.Name = file.Filename
			infoFile.Size = file.Size
			infoFile.Status = err.Error()
			resp = append(resp, infoFile)
			break
		}
		infoFile.Id = idFile
		infoFile.Name = file.Filename
		infoFile.Size = file.Size
		infoFile.Status = "success"
		resp = append(resp, infoFile)
	}
	ctx.JSON(http.StatusOK, resp)
}
