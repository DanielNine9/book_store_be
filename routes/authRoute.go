package routes

import (
	"shop-account/handlers"
	"shop-account/middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	authGroup := router.Group("/auth")
	{
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	authGroup.PUT("/update-role", middlewares.AuthMiddlewareForRole("admin"), authHandler.UpdateRole)
	}
}
