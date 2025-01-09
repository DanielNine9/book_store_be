package routes

import (
    "shop-account/handlers"
    "shop-account/middlewares"
    "github.com/gin-gonic/gin"
)

func TransactionRoutes(router *gin.Engine, transactionHandler *handlers.TransactionHandler) {
    transactionGroup := router.Group("/transactions")
    transactionGroup.Use(middlewares.AuthMiddleware())
    {
        transactionGroup.POST("/", transactionHandler.CreateTransaction)          
        transactionGroup.GET("/", transactionHandler.GetUserTransactions)        
        // transactionGroup.GET("/:id", purchaseHandler.GetTransactionByID)       
        // transactionGroup.PUT("/:id", purchaseHandler.UpdateTransactionStatus) 
        transactionGroup.DELETE("/:id", transactionHandler.DeleteTransaction)   
    }
}
