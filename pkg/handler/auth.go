package handler

import (
	cloud_storage "github.com/danil-vas/cloud-storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary SingUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body cloud_storage.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /auth/sing-up [post]
func (h *Handler) singUp(c *gin.Context) {
	var input cloud_storage.User

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if len(input.Login) < 4 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "minimum login length 4",
		})
		return
	}
	if len(input.Password) < 4 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "minimum password length 4",
		})
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	err = h.services.Authorization.CreateMainDirectory(id, input.Login)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Login    string `json:"login" building:"required"`
	Password string `json:"password" building:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {string}  string    "error"
// @Failure 500 {string}  string    "error"
// @Failure default {string}  string    "error"
// @Router /auth/sign-in [post]
func (h *Handler) singIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Login, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "login or password entered incorrectly",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
