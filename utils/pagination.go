package utils

import (
	"fmt"
	"math"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

// PaginateAndSearch handles pagination and search for any model
func PaginateAndSearch(c *gin.Context, db *gorm.DB, model interface{}, result interface{}) (int64, int, int, error) {
	// Get pagination parameters from query
	pageStr := c.DefaultQuery("page", "1")  // Default to page 1 if not provided
	limitStr := c.DefaultQuery("limit", "10") // Default to limit 10 if not provided
	search := c.DefaultQuery("search", "")  // Default to empty string if no search query provided

	// Parse page and limit to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return 0, 0, 0, fmt.Errorf("invalid page number")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 0, 0, 0, fmt.Errorf("invalid limit")
	}

	// Calculate the offset for pagination
	offset := (page - 1) * limit

	// Initialize the query builder
	query := db.Model(model)

	// Apply search filter if provided (search by name and description for example)
	if search != "" {
		search = "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	// Get the total number of records (for pagination)
	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return 0, 0, 0, err
	}

	// Fetch the records into the provided result slice
	if err := query.Limit(limit).Offset(offset).Find(result).Error; err != nil {
		return 0, 0, 0, err
	}

	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	return totalItems, page, totalPages, nil
}
