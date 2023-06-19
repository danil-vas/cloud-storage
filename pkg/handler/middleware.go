package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "empty auth header",
		})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid auth header",
		})
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])

	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.Set(userCtx, userId)
}

func getUsersId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id not found")
	}
	return idInt, nil
}
