package middlewares

import (
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

func SetUserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString := strings.Split(authHeader, " ")[1]

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte("your_secret_key"), nil
			})

			if err == nil && token.Valid {
				if userID, ok := claims["user_id"].(float64); ok {
					c.Set("id", uint(userID))
				}

			}
		}
		c.Next()
	}
}
