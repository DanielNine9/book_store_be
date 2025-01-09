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
		purchaseGroup.GET("/", purchaseHandler.GetUserPurchases) 
		purchaseGroup.DELETE("/:purchase_id", purchaseHandler.DeletePurchase)
		purchaseGroup.PUT("/:purchase_id", purchaseHandler.UpdatePurchase)
	}
}
