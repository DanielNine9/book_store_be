package routes

import (
	"shop-account/handlers"
	"github.com/gin-gonic/gin"
)

// AuthRoutes đăng ký các route liên quan đến xác thực
func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	// Đăng ký route cho đăng nhập
	router.POST("/login", authHandler.Login)

	// Đăng ký route cho đăng ký người dùng (nếu có)
	// router.POST("/register", authHandler.Register)
}
