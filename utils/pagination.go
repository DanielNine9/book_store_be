package utils

import (
	"fmt"
	"math"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

// PaginateAndSearch handles pagination, dynamic search on custom fields, and custom query operators
func PaginateAndSearch(c *gin.Context, db *gorm.DB, model interface{}, result interface{}, customQuery *gorm.DB) (int64, int, int, error) {
	// Get pagination parameters from query
	pageStr := c.DefaultQuery("page", "1")  // Default to page 1 if not provided
	limitStr := c.DefaultQuery("limit", "10") // Default to limit 10 if not provided
	search := c.DefaultQuery("search", "")  // Default to empty string if no search query provided
	searchFields := c.DefaultQuery("search_fields", "") // Fields to search by (comma-separated)
	searchOperator := c.DefaultQuery("search_operator", "OR") // Operator to combine search conditions

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

	// Apply the custom query if provided
	if customQuery != nil {
		query = customQuery
	}

	// Apply dynamic search if search is provided
	if search != "" && searchFields != "" {
		// Split search fields into a list
		fields := strings.Split(searchFields, ",")
		operator := strings.ToUpper(searchOperator)
		if operator != "AND" && operator != "OR" {
			operator = "OR" // Default to OR if an invalid operator is provided
		}

		// Build the search conditions dynamically
		var conditions []string
		var args []interface{}
		for _, field := range fields {
			conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE ?", field))
			args = append(args, "%"+strings.ToLower(search)+"%")
		}

		// Combine the conditions using the specified operator (AND/OR)
		conditionString := strings.Join(conditions, fmt.Sprintf(" %s ", operator))

		// Apply the dynamic search condition
		query = query.Where(conditionString, args...)
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
