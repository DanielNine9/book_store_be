package admin

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"shop-account/models"
	"shop-account/utils"
)

type AdminTransactionHandler struct {
	DB *gorm.DB
}

// GetAllTransactions allows admins to view all transactions
func (h *AdminTransactionHandler) GetAllTransactions(c *gin.Context) {
	var transactions []models.Transaction

	// Get pagination and search parameters from the query string
	totalItems, page, totalPages, err := utils.PaginateAndSearch(c, h.DB, &models.Transaction{}, &transactions, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions", "details": err.Error()})
		return
	}

	// Return the paginated transactions along with metadata
	c.JSON(http.StatusOK, gin.H{
		"current_page":   page,
		"total_pages":    totalPages,
		"total_items":    totalItems,
		"items_per_page": c.DefaultQuery("limit", "10"),
		"transactions":   transactions,
	})
}

func (h *AdminTransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	}

	var requestBody struct {
		Status models.TransactionStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

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
