package routes

import (
	"shop-account/handlers/admin"
	"shop-account/middlewares"
	"github.com/gin-gonic/gin"

)

func AdminRoutes(router *gin.Engine, adminTransactionHandler *admin.AdminTransactionHandler) {
	adminGroup := router.Group("/admin")
    adminGroup.Use(middlewares.AuthMiddlewareForRole("admin"))

	{
		adminGroup.GET("/transactions", adminTransactionHandler.GetAllTransactions)
		adminGroup.PATCH("/transactions/:id/status", adminTransactionHandler.UpdateTransactionStatus)
	}
}

