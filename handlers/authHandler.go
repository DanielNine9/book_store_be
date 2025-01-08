package handlers

import (
	"shop-account/models"
	"shop-account/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler struct cho việc đăng nhập
type AuthHandler struct {
	// Nếu bạn có một DB user, có thể thêm tham số DB ở đây
}

// Hàm đăng nhập, tạo token
func (h *AuthHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Kiểm tra user (ví dụ, so sánh với dữ liệu giả hoặc trong database)
	// Lưu ý: Mật khẩu cần được hash trong thực tế
	if user.Username != "admin" || user.Password != "password" { // Sử dụng thông tin giả
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Sinh JWT token
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Trả về token cho người dùng
	c.JSON(http.StatusOK, gin.H{"token": token})
}
