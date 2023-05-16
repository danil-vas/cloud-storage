package handler

import "github.com/gin-gonic/gin"

func (h *Handler) downloadFile(c *gin.Context) {
	c.Set("Content-Disposition", "attachment; filename=index.html")
	c.Set("Content-Type", "application/octet-stream")
	c.FileAttachment("temp/index.html", "index.html")
}
