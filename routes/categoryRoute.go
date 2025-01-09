package routes

import (
	"shop-account/handlers"
	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine, categoryHandler *handlers.CategoryHandler) {
	categoryRoutes := router.Group("/categories")
	{
		categoryRoutes.POST("/", categoryHandler.CreateCategory)
		categoryRoutes.GET("/", categoryHandler.GetCategories)
		categoryRoutes.GET("/:id", categoryHandler.GetCategory)
		categoryRoutes.PUT("/:id", categoryHandler.UpdateCategory)
		categoryRoutes.PATCH("/:id", categoryHandler.PatchCategory)
		categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
