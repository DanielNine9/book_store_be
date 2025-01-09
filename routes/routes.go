package routes

import (
	"shop-account/handlers"
	"shop-account/handlers/admin"
	"github.com/gin-gonic/gin"
)

// SetupRoutes đăng ký tất cả các route cho API, bao gồm cả xác thực
func SetupRoutes(router *gin.Engine,categoryHandler *handlers.CategoryHandler,adminTransactionHandler *admin.AdminTransactionHandler, transactionHandler *handlers.TransactionHandler, purchaseHandler *handlers.PurchaseHandler, userHandler *handlers.UserHandler, authorHandler *handlers.AuthorHandler, bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler) {
	AuthorRoutes(router, authorHandler)

	BookRoutes(router, bookHandler)
	CategoryRoutes(router, categoryHandler)
	AuthRoutes(router, authHandler)
	UserRoutes(router, userHandler)
	PurchaseRoutes(router, purchaseHandler)
	TransactionRoutes(router, transactionHandler)
	AdminRoutes(router, adminTransactionHandler)
}
