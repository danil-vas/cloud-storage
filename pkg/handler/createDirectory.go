package handler

import (
	"fmt"
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func (h *Handler) createDirectory(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	path = "temp/" + path
	nameDirectory, _ := c.GetPostForm("name")
	err = os.MkdirAll(path+"/"+nameDirectory, 0777)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.Directory.AddDirectory(userId, objectId, nameDirectory)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "complete",
	})
}

// @Summary Get Directory
// @Security ApiKeyAuth
// @Tags Directory
// @Description get directory
// @ID get-directory
// @Accept  json
// @Produce  json
// @Success 200 {object} cloud_storage.Node
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /api/directory/{id} [get]
func (h *Handler) getDirectory(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
	resp, err := h.services.Directory.GetDirectoriesAndFiles(userId, objectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Get Main Directory
// @Security ApiKeyAuth
// @Tags Directory
// @Description get main directory
// @ID get-main-directory
// @Accept  json
// @Produce  json
// @Success 200 {object} cloud_storage.Node
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /api/directory [get]
func (h *Handler) getMainDirectory(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	objectId, err := h.services.Directory.GetIdMainDirectory(userId)
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
	resp, err := h.services.Directory.GetDirectoriesAndFiles(userId, objectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

func traverseNodes(node *cloud_storage.Node, user_id int, h *Handler) {
	if len(node.Children) == 0 {
		// Добавьте здесь нужную обработку самого дальнего элемента
		if node.Type == "file" {
			path, err := h.services.File.PathUploadFile(user_id, node.ID)
			fmt.Println(path)
			if err != nil {
				return
			}
			path = "temp/" + path
			err = os.Remove(path)
			if err != nil {
				return
			}
			err = h.services.File.DeleteFile(user_id, node.ID)
		} else {
			path, err := h.services.File.PathUploadFile(user_id, node.ID)
			if err != nil {
				return
			}
			path = "temp/" + path
			err = os.Remove(path)
			if err != nil {
				return
			}
			err = h.services.Directory.DeleteDirectory(node.ID)
		}
		return
	}

	for _, child := range node.Children {
		traverseNodes(child, user_id, h)
	}

	// Добавьте здесь нужную обработку родительской директории
	if node.Type == "file" {
		path, err := h.services.File.PathUploadFile(user_id, node.ID)
		if err != nil {
			return
		}
		path = "temp/" + path
		err = os.Remove(path)
		if err != nil {
			return
		}
		err = h.services.File.DeleteFile(user_id, node.ID)
	} else {
		path, err := h.services.File.PathUploadFile(user_id, node.ID)
		if err != nil {
			return
		}
		path = "temp/" + path
		err = os.Remove(path)
		if err != nil {
			return
		}
		err = h.services.Directory.DeleteDirectory(node.ID)
	}
}

func (h *Handler) deleteDirectory(c *gin.Context) {
	userId, err := getUsersId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
	resp, err := h.services.Directory.GetDirectoriesAndFiles(userId, objectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, item := range resp {
		traverseNodes(&item, userId, h)
	}
	path, err := h.services.File.PathUploadFile(userId, objectId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	path = "temp/" + path
	err = os.Remove(path)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	err = h.services.File.DeleteFile(userId, objectId)
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "complete",
	})
}
