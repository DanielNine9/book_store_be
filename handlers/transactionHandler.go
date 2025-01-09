package handlers

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"shop-account/models"
	"github.com/lib/pq"
)

type TransactionHandler struct {
	DB *gorm.DB
}
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
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

	var purchaseRequest struct {
		PurchaseIDs []uint `json:"purchase_ids"` 
	}

	if err := c.ShouldBindJSON(&purchaseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Printf("purchaseRequest PurchaseIDs: %v\n", purchaseRequest.PurchaseIDs)

	var purchases []models.Purchase
	if err := h.DB.Preload("Book").Where("id = ANY(?) and transaction_id = 0", pq.Array(purchaseRequest.PurchaseIDs)).Find(&purchases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve purchases"})
		return
	}

	if len(purchases) != len(purchaseRequest.PurchaseIDs) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Some purchases not found"})
		return
	}

	var totalAmount float64
	for _, purchase := range purchases {
		
	fmt.Printf("Quantity: %f\n", float64(purchase.Quantity))
	fmt.Printf("BookPrice: %f\n", float64(purchase.BookPrice))
		totalAmount += float64(purchase.Quantity) * float64(purchase.BookPrice) 
	}

	fmt.Printf("totalAmount: %f\n", totalAmount)

	transaction := models.Transaction{
		UserID:      userID,           
		TotalAmount: totalAmount,     
		Status:      "pending",        
		Purchases:   purchases,        
	}

	if err := h.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a new transaction"})
		return
	}

	for i := range purchases {
        purchases[i].TransactionID = transaction.ID
        if err := h.DB.Save(&purchases[i]).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase with transaction ID"})
            return
        }
    }

	c.JSON(http.StatusOK, gin.H{
		"message":    "A transaction is created",
	})
}
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
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

    transactionID := c.Param("id")
    if transactionID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
        return
    }

    var transaction models.Transaction
    if err := h.DB.Preload("Purchases").Preload("Purchases.Book").First(&transaction, transactionID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
        return
    }

    if transaction.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this transaction"})
        return
    }

    for _, purchase := range transaction.Purchases {
        purchase.TransactionID = 0 
        if err := h.DB.Save(&purchase).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase records"})
            return
        }
    }

    if err := h.DB.Delete(&transaction).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}



func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
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

	var transactions []models.Transaction
	if err := h.DB.Preload("Purchases").Preload("Purchases.Book").Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	if len(transactions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No transactions found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}
