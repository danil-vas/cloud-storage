package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Get User
// @Security ApiKeyAuth
// @Tags User
// @Description get user
// @ID get-user
// @Accept  json
// @Produce  json
// @Success 200 {object} cloud_storage.UserInfo
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /api/user/ [get]
func (h *Handler) getUser(ctx *gin.Context) {
	userId, err := getUsersId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error authorization"})
		return
	}
	info, err := h.services.User.GetUser(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, info)
}
