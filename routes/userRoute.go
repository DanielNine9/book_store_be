package routes

import (
	"shop-account/handlers" 
	"github.com/gin-gonic/gin"
)

// UserRoutes registers routes for the User resource
func UserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	userGroup := router.Group("/users")
	// You can add authentication middleware if needed, e.g., userGroup.Use(middlewares.AuthMiddleware())
	{
		userGroup.GET("/", userHandler.List)          
		userGroup.GET("/:id", userHandler.GetByID)    
		userGroup.POST("/", userHandler.Create)       
		userGroup.PUT("/:id", userHandler.Update)     
		userGroup.DELETE("/:id", userHandler.Delete)  
	}
}
