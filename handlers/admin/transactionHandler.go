package admin

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"shop-account/models"
)

type AdminTransactionHandler struct {
	DB *gorm.DB
}

// GetAllTransactions allows admins to view all transactions
func (h *AdminTransactionHandler) GetAllTransactions(c *gin.Context) {
	var transactions []models.Transaction

	// Retrieve all transactions with related purchases and books
	if err := h.DB.Preload("Purchases").Preload("Purchases.Book").Find(&transactions).Error; err != nil {
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

// UpdateTransactionStatus allows admins to change the status of a transaction
func (h *AdminTransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	// Get transaction ID from URL parameters
	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	// Get the new status from the request body
	var requestBody struct {
		Status models.TransactionStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate that the status is valid
	validStatuses := []models.TransactionStatus{
		models.Pending,
		models.Approved,
		models.Rejected,
		models.Completed,
	}

	isValidStatus := false
	for _, validStatus := range validStatuses {
		if requestBody.Status == validStatus {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Find the transaction by ID
	var transaction models.Transaction
	if err := h.DB.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Update the status of the transaction
	transaction.Status = requestBody.Status
	if err := h.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Transaction status updated successfully",
		"transaction": transaction,
	})
}
