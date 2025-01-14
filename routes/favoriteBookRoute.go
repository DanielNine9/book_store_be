package routes

import (
    "shop-account/handlers"
    "shop-account/middlewares"
    "github.com/gin-gonic/gin"
)

func FavoriteBookRoutes(router *gin.Engine, favoriteBookHandler *handlers.FavoriteBookHandler) {
    favoriteBookGroup := router.Group("/favorites")
    favoriteBookGroup.Use(middlewares.AuthMiddleware())
    {
        favoriteBookGroup.POST("/", favoriteBookHandler.CreateFavoriteBook) 
        favoriteBookGroup.GET("/", favoriteBookHandler.GetUserFavoriteBooks) 
        favoriteBookGroup.DELETE("/:id", favoriteBookHandler.DeleteFavoriteBook) 
    }
}
