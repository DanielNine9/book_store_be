package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"shop-account/models"
	"strconv"
	"fmt"
)

// BookHandler chứa các phương thức xử lý cho sách.
type PurchaseHandler struct {
    DB *gorm.DB
}
func (h *PurchaseHandler) BuyBook(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    bookID := c.Param("book_id")
	id, err := strconv.Atoi(bookID)


	if err != nil {
		fmt.Printf("Error converting bookID: %s, Error: %v\n", bookID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Error "})
        return
	}
    var book models.Book
    if err := h.DB.Where("id = ? AND active = ?", id, true).First(&book).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or inactive"})
        return
    }
	
	fmt.Printf("user_id: %s\n", userID)
    var user models.User
    if err := h.DB.Where("id = ? AND active = ?", userID, true).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found or inactive"})
        return
    }

    // Tạo bản ghi mua sách
    purchase := models.Purchase{
        UserID: user.ID,
        BookID: book.ID,
        User:   user,
        Book:   book,
    }

    if err := h.DB.Create(&purchase).Error; err != nil {
		
	fmt.Printf("err: %s\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purchase book"})
        return
    }

    // Trả về phản hồi thành công
    c.JSON(http.StatusOK, gin.H{
        "message": "Book purchased successfully",
        "purchase": purchase,
    })
}