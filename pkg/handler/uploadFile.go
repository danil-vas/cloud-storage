package handler

import (
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

func (h *Handler) uploadFile(ctx *gin.Context) {
	userId, err := getUsersId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error authorization"})
		return
	}
	objectId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	flag, err := h.services.CheckAccessToObject(userId, objectId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if flag == false {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "Forbidden"})
		return
	}
	path, err := h.services.File.PathUploadFile(userId, objectId)
	form, _ := ctx.MultipartForm()
	files := form.File["myFile"]
	str := make([]string, 0)
	resp := make([]jsonResponse, 0)
	for _, file := range files {
		var infoFile jsonResponse
		size := file.Size
		availableMemory, err := h.services.File.GetAvailableMemory(userId)
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
