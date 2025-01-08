package middlewares

import (
	"shop-account/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
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


func AuthMiddlewareForRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Split the header into "Bearer token"
		tokenString := strings.Split(authHeader, " ")[1]

		// Parse and validate the JWT token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the role matches
		if claims["role"] != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Proceed if role matches
		c.Next()
	}
}
