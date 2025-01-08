package routes

import (
	"shop-account/handlers"
	// "shop-account/middlewares"
	"github.com/gin-gonic/gin"
)

// AuthorRoutes đăng ký các route cho Author
func AuthorRoutes(router *gin.Engine, authorHandler *handlers.AuthorHandler) {
	authorGroup := router.Group("/authors")
	// authorGroup.Use(middlewares.AuthMiddleware())
	{
		authorGroup.GET("/", authorHandler.GetAuthors)
		authorGroup.GET("/:id", authorHandler.GetAuthorByID)
		authorGroup.POST("/", authorHandler.CreateAuthor)
		authorGroup.PUT("/:id", authorHandler.UpdateAuthor)
		authorGroup.PATCH("/:id", authorHandler.PatchAuthor) // Thêm PATCH nếu cần
		authorGroup.DELETE("/:id", authorHandler.DeleteAuthor)
	}
}
