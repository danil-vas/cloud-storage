package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func (h *Handler) deleteFile(c *gin.Context) {
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
	path, err := h.services.File.PathUploadFile(userId, objectId)
	if err != nil {
		return
	}
	path = "temp/" + path
	err = os.Remove(path)
	if err != nil {
		return
	}
	err = h.services.File.DeleteFile(userId, objectId)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "complete",
	})
}
