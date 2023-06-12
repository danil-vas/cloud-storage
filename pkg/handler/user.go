package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getUser(ctx *gin.Context) {
	userId, err := getUsersId(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error authorization"})
		return
	}
	info, err := h.services.User.GetUser(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, info)
}
