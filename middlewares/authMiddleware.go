package middlewares

import (
	"shop-account/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Middleware để kiểm tra JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Tách "Bearer" khỏi token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Kiểm tra và xác thực token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Đặt thông tin người dùng vào context để dùng ở các route sau
		c.Set("username", claims.Username)
		c.Next()
	}
}
