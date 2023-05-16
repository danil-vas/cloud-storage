package handler

import (
	"github.com/danil-vas/cloud-storage/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sing-up", h.singUp)
		auth.POST("/sing-in", h.singIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		file := api.Group("/file")
		{
			file.POST("/:id", h.uploadFile)
		}
		//api.POST("/upload", h.uploadFile)
		api.GET("/download", h.downloadFile)
	}

	return router
}
