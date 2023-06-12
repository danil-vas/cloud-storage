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
)

func (h *Handler) downloadFile(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		return
	}
	objectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	flag, err := h.services.CheckAccessToObject(userId, objectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if flag == false {
		c.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}
	typeFile, err := h.services.GetTypeObject(objectId)
	path, err := h.services.File.PathUploadFile(userId, objectId)
	if err != nil {
		return
	}
	path = "temp/" + path
	if typeFile == "file" {
		file, err := os.Open(path)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		originalFileName, err := h.services.File.OriginalFileName(objectId)
		contentDisposition := fmt.Sprintf("attachment; filename=%s", originalFileName)
		c.Header("Content-Disposition", contentDisposition)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.File(path)
	} else {
		zipFilePath := path + "Archive.zip"
		zipFile, err := os.Create(zipFilePath)
		if err != nil {
			panic(err)
		}
		defer zipFile.Close()

		archive := zip.NewWriter(zipFile)
		defer archive.Close()
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
			zipFile, err := archive.Create(relPath)
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
			panic(err)
		}
		file, err := os.Open(zipFilePath)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		originalFileName, err := h.services.File.OriginalFileName(objectId)
		originalFileName += "Archive.zip"
		contentDisposition := fmt.Sprintf("attachment; filename=%s", originalFileName)
		c.Header("Content-Disposition", contentDisposition)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	}

}
