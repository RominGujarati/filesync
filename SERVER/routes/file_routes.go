package routes

import (
	"server/file"
	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.Engine) {
	fileGroup := router.Group("/file")
	{
		fileGroup.GET("/", file.GetAllFiles)
		fileGroup.GET("/:fileName", file.GetFileContent)
	}
}
