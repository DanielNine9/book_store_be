package routes

import (
	"shop-account/handlers"
	"shop-account/middlewares"
	"github.com/gin-gonic/gin"
)

// BookRoutes đăng ký các route cho sách
func BookRoutes(router *gin.Engine, bookHandler *handlers.BookHandler) {
	bookGroup := router.Group("/books")
	bookGroup.Use(middlewares.SetUserIDMiddleware())
	{
		bookGroup.GET("/", bookHandler.GetBooks)
		bookGroup.GET("/:id", bookHandler.GetBookByID)
		bookGroup.POST("/", bookHandler.CreateBook)
		bookGroup.PUT("/:id", bookHandler.UpdateBook)
		bookGroup.PUT("/restore/:id", bookHandler.Restore)
		bookGroup.PATCH("/:id", bookHandler.PatchBook)
		bookGroup.DELETE("/:id", bookHandler.DeleteBook)
		bookGroup.GET("/concurrency", bookHandler.GetBooksConcurrently) 
		bookGroup.GET("/not-concurrency", bookHandler.GetBooksNotConcurrently) 
		bookGroup.GET("/import", bookHandler.ImportBooks) 
	}
}
