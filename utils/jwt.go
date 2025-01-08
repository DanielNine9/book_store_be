package utils

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	// "log"
)

var SecretKey = []byte("your-secret-key") // Đặt một secret key cho JWT

// Struct cho claims của JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Hàm tạo JWT token
func GenerateToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token có hiệu lực trong 24 giờ
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// Hàm kiểm tra và parse JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
