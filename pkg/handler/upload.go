package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
)

func (h *Handler) uploadFile(ctx *gin.Context) {
	userId, err := getUsersId(ctx)
	if err != nil {
		return
	}
	objectId, err := strconv.Atoi(ctx.Param("id"))
	path, err := h.services.File.PathUploadFile(userId, objectId)
	form, _ := ctx.MultipartForm()
	files := form.File["myFile"]
	str := make([]string, 0)
	for _, file := range files {
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
			fmt.Println(err)
			return
		}
		defer tempFile.Close()
		readerFile, _ := file.Open()
		_, err = io.Copy(tempFile, readerFile)
		if err != nil {
			fmt.Println(err)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"filepath": path})
}
