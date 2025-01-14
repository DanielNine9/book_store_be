package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"shop-account/models"
	"strconv"
	"fmt"
)

type FavoriteBookHandler struct {
	DB *gorm.DB
}

func (h *FavoriteBookHandler) CreateFavoriteBook(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	fmt.Printf("userIDInterface %s \n", userIDInterface)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	userID := uint(userIDFloat)

	var request struct {
		BookID uint `json:"book_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var book models.Book
	if err := h.DB.First(&book, request.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var favorite models.FavoriteBook
	if err := h.DB.Where("user_id = ? AND book_id = ?", userID, request.BookID).First(&favorite).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book already in favorites"})
		return
	}

	favoriteBook := models.FavoriteBook{
		UserID: userID,
		BookID: request.BookID,
	}

	if err := h.DB.Create(&favoriteBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Book added to favorites",
		"favorite": favoriteBook.ID, // Use the newly created favoriteBook ID
	})
}
func (h *FavoriteBookHandler) DeleteFavoriteBook(c *gin.Context) {
	// Get the authenticated user ID from the context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Ensure the user ID is of type uint
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	userID := uint(userIDFloat)

	// Get the favorite ID from the route parameter
	favoriteIDStr := c.Param("id")
	if favoriteIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Favorite ID is required"})
		return
	}

	// Convert favoriteID from string to uint
	favoriteID, err := strconv.Atoi(favoriteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid favorite ID"})
		return
	}
	// Find the favorite book record by favorite_id and user_id
	var favoriteBook models.FavoriteBook
	if err := h.DB.Where("id = ? AND user_id = ?", favoriteID, userID).First(&favoriteBook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite book not found"})
		return
	}

	// Delete the favorite book record
	if err := h.DB.Delete(&favoriteBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove book from favorites"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Book removed from favorites",
	})
}


func (h *FavoriteBookHandler) GetUserFavoriteBooks(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}
	userID := uint(userIDFloat)

	var favoriteBooks []models.FavoriteBook
	if err := h.DB.Preload("Book").Where("user_id = ?", userID).Find(&favoriteBooks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve favorite books"})
		return
	}

	if len(favoriteBooks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No favorite books found"})
		return
	}

	var books []models.Book
	for _, favorite := range favoriteBooks {
		books = append(books, favorite.Book)
	}

	c.JSON(http.StatusOK, gin.H{
		"favorite_books": favoriteBooks,
	})
}
