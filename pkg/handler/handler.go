package handler

import (
	_ "github.com/danil-vas/cloud-storage/docs"
	"github.com/danil-vas/cloud-storage/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
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
			file.GET("/:id", h.downloadFile)
			file.DELETE("/:id", h.deleteFile)
		}
		directory := api.Group("/directory")
		{
			directory.POST("/:id", h.createDirectory)
			directory.GET("/:id", h.getDirectory)
			directory.GET("/", h.getMainDirectory)
			directory.DELETE("/:id", h.deleteDirectory)
		}
		user := api.Group("/user")
		{
			user.GET("/", h.getUser)
			user.POST("/share/:id", h.shareObject)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
