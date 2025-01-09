package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	// "fmt"
)

// Middleware để kiểm tra JWT token
	func AuthMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
				c.Abort()
				return
			}

			tokenString := strings.Split(authHeader, " ")[1]

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte("your_secret_key"), nil
			})

			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
			// fmt.Printf("claims: %s", claims)

			c.Set("username", claims["username"])
			c.Set("user_id", claims["user_id"])
			c.Next()
		}
	}


func AuthMiddlewareForRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if claims["role"] != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
