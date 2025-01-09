package routes

import (
	"shop-account/handlers"
	"github.com/gin-gonic/gin"
	"shop-account/middlewares"
)

func PurchaseRoutes(router *gin.Engine, purchaseHandler *handlers.PurchaseHandler) {
	purchaseGroup := router.Group("/purchases")
	purchaseGroup.Use(middlewares.AuthMiddleware())
	{
		purchaseGroup.POST("/:book_id", purchaseHandler.BuyBook) 
	}
}
