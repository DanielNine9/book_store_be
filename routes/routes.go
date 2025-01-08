package routes

import (
	"shop-account/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes đăng ký tất cả các route cho API, bao gồm cả xác thực
func SetupRoutes(router *gin.Engine, authorHandler *handlers.AuthorHandler, bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler) {
	// Đăng ký các route cho Author
	AuthorRoutes(router, authorHandler)

	// Đăng ký các route cho Book
	BookRoutes(router, bookHandler)

	// Đăng ký các route cho Auth (Xác thực)
	AuthRoutes(router, authHandler)
}
