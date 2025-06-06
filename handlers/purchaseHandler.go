package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"shop-account/models"
	"shop-account/utils"
	"shop-account/dtos"
	"strconv"
	"fmt"
	"time"
)

type PurchaseHandler struct {
    DB *gorm.DB
}
func (h *PurchaseHandler) BuyBook(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var purchaseRequest dtos.PurchaseRequest
    if err := c.ShouldBindJSON(&purchaseRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    bookID := c.Param("book_id")
    id, err := strconv.Atoi(bookID)
    if err != nil {
        fmt.Printf("Error converting bookID: %s, Error: %v\n", bookID, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Error converting book ID"})
        return
    }

    var book models.Book
    if err := h.DB.Where("id = ? AND active = ?", id, true).First(&book).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or inactive"})
        return
    }

    if book.QuantityInStock < purchaseRequest.Quantity {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": fmt.Sprintf("Not enough stock available, the quantity that can be chosen is %d", book.QuantityInStock),
        })
        return
    }

    var user models.User
    if err := h.DB.Where("id = ? AND active = ?", userID, true).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found or inactive"})
        return
    }

    purchase := models.Purchase{
        UserID:   user.ID,
        BookID:   book.ID,
        Quantity: purchaseRequest.Quantity,
        BookPrice: book.Price,
    }

    // Generate unique code for the purchase
    code, err := utils.GenerateCode(h.DB, &models.Purchase{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate purchase code", "details": err.Error()})
        return
    }
    purchase.Code = code

    if err := h.DB.Create(&purchase).Error; err != nil {
        fmt.Printf("Error creating purchase: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purchase book"})
        return
    }

    book.QuantityInStock -= purchaseRequest.Quantity
    book.QuantitySold += purchaseRequest.Quantity

    if err := h.DB.Save(&book).Error; err != nil {
        fmt.Printf("Error updating book: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book stock"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Book purchased successfully",
        "purchase": purchase,
    })
}


// func (h *PurchaseHandler) BuyBook(c *gin.Context) {
//     userID, exists := c.Get("user_id")
//     if !exists {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
//         return
//     }
// 	var purchaseRequest dtos.PurchaseRequest
//     if err := c.ShouldBindJSON(&purchaseRequest); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//         return
//     }

//     bookID := c.Param("book_id")
// 	id, err := strconv.Atoi(bookID)


// 	if err != nil {
// 		fmt.Printf("Error converting bookID: %s, Error: %v\n", bookID, err)
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Error "})
//         return
// 	}
//     var book models.Book
//     if err := h.DB.Where("id = ? AND active = ?", id, true).First(&book).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or inactive"})
//         return
//     }
	
// 	fmt.Printf("user_id: %s\n", userID)
//     var user models.User
//     if err := h.DB.Where("id = ? AND active = ?", userID, true).First(&user).Error; err != nil {
//         c.JSON(http.StatusNotFound, gin.H{"error": "User not found or inactive"})
//         return
//     }

//     purchase := models.Purchase{
//         UserID: user.ID,
//         BookID: book.ID,
//         User:   user,
//         Book:   book,
// 		Quantity: purchaseRequest.Quantity,
//     }

//     if err := h.DB.Create(&purchase).Error; err != nil {
		
// 	fmt.Printf("err: %s\n", err)
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purchase book"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{
//         "message": "Book purchased successfully",
//         "purchase": purchase,
//     })
// }

func (h *PurchaseHandler) GetUserPurchases(c *gin.Context) {
    // Get the authenticated user ID from context
    userIDInterface, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Convert the user ID to uint
    userIDFloat, ok := userIDInterface.(float64)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
        return
    }
    userID := uint(userIDFloat)

    // Initialize an empty slice to hold purchases
    var purchases []models.Purchase

    // Preload related data and filter by user_id
    query := h.DB.Preload("Book").Preload("User").Preload("Transaction").Where("user_id = ?", userID)

    // Call PaginateAndSearch utility to fetch paginated data with dynamic search (if any)
    totalItems, page, totalPages, err := utils.PaginateAndSearch(c, query, &models.Purchase{}, &purchases, nil)
    if err != nil {
        fmt.Printf("Error fetching purchases for user_id %v: %v\n", userID, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch purchases"})
        return
    }

    // If no purchases found, return empty response
    if len(purchases) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No purchases found for this user", "purchases": []dtos.PurchaseResponse{}})
        return
    }

    // Convert the purchases into PurchaseResponse DTOs
    var purchaseResponses []dtos.PurchaseResponse
    for _, purchase := range purchases {
        var deletedAt *string
        if purchase.DeletedAt != nil {
            deletedAtStr := purchase.DeletedAt.Format(time.RFC3339)
            deletedAt = &deletedAtStr
        }

        purchaseResponses = append(purchaseResponses, dtos.PurchaseResponse{
            ID:        purchase.ID,
            CreatedAt: purchase.CreatedAt.Format(time.RFC3339),
            UpdatedAt: purchase.UpdatedAt.Format(time.RFC3339),
            DeletedAt: deletedAt,
            UserID:    purchase.UserID,
            Book:      purchase.Book,
            Quantity:  purchase.Quantity,
        })
    }

    // Return the paginated purchases response
    c.JSON(http.StatusOK, gin.H{
        "message":   "User purchases retrieved successfully",
        "current_page":   page,
        "total_pages":    totalPages,
        "total_items":    totalItems,
        "items_per_page": c.DefaultQuery("limit", "10"),
        "purchases":      purchaseResponses,
    })
}


func (h *PurchaseHandler) UpdatePurchase(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    purchaseID := c.Param("purchase_id")
    id, err := strconv.Atoi(purchaseID)
    if err != nil {
        fmt.Printf("Error converting purchaseID: %s, Error: %v\n", purchaseID, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Invalid purchase ID"})
        return
    }

    var purchase models.Purchase
    if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&purchase).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Purchase not found or you don't have access"})
        return
    }

    var input struct {
        Quantity uint `json:"quantity"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    purchase.Quantity = input.Quantity
    if err := h.DB.Save(&purchase).Error; err != nil {
        fmt.Printf("Error updating purchase: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "Purchase updated successfully",
        "purchase": purchase,
    })
}

// DeletePurchase - Marks a purchase as deleted (soft delete)
func (h *PurchaseHandler) DeletePurchase(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    purchaseID := c.Param("purchase_id")
    id, err := strconv.Atoi(purchaseID)
    if err != nil {
        fmt.Printf("Error converting purchaseID: %s, Error: %v\n", purchaseID, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Invalid purchase ID"})
        return
    }

    var purchase models.Purchase
    if err := h.DB.Where("id = ? AND user_id = ?", id, userID).First(&purchase).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Purchase not found or you don't have access"})
        return
    }

    // Perform soft delete by setting the DeletedAt timestamp
    if err := h.DB.Model(&purchase).Update("deleted_at", time.Now()).Error; err != nil {
        fmt.Printf("Error deleting purchase: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete purchase"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Purchase deleted successfully",
    })
}