package routes

import (
	"server/auth"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
	}
}

func ProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/protected")
	protected.Use(auth.JWTAuthMiddleware()) // Apply JWT middleware
	{
		protected.GET("/profile", auth.Profile)
		protected.POST("/upload-yaml", auth.UploadYAML)
	}
}
