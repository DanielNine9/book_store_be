package handlers

import (
	"shop-account/models"
	"shop-account/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
	"fmt"
)

// Handler struct để lưu trữ db instance
type AuthorHandler struct {
	DB *gorm.DB
}

// Hàm lấy danh sách tác giả
func (h *AuthorHandler) GetAuthors(c *gin.Context) {
	// Initialize an empty slice to hold the authors
	var authors []models.Author

	// Preload related data (Books)
	query := h.DB.Preload("Books")

	// Call PaginateAndSearch utility to fetch paginated data with dynamic search (if any)
	totalItems, page, totalPages, err := utils.PaginateAndSearch(c, query, &models.Author{}, &authors, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch authors", "details": err.Error()})
		return
	}

	// Returning the paginated authors as a response
	c.JSON(http.StatusOK, gin.H{
		"current_page":   page,
		"total_pages":    totalPages,
		"total_items":    totalItems,
		"items_per_page": c.DefaultQuery("limit", "10"),
		"authors":        authors,
	})
}


// Hàm tạo tác giả mới
func (h *AuthorHandler) CreateAuthor(c *gin.Context) {
	var author models.Author
	// Bind dữ liệu JSON từ request body
	if err := c.ShouldBindJSON(&author); err != nil {
		// Log the detailed error
		fmt.Printf("Error binding JSON: %s\n", err.Error())
		// Return detailed error message to client (for development purposes)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error binding JSON: %s", err.Error())})
		return
	}

	// Kiểm tra nếu name bị bỏ trống
	if strings.TrimSpace(author.Name) == "" {
		errMsg := "Author name is required"
		// Log the detailed error
		fmt.Printf("%s\n", errMsg)
		// Return the detailed error to the client
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	// Generate unique code for the author
	code, err := utils.GenerateCode(h.DB, &models.Author{})
	if err != nil {
		// Log the detailed error from GenerateCode
		fmt.Printf("Error generating author code: %s\n", err.Error())
		// Return detailed error message to client (for development purposes)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate author code: %s", err.Error())})
		return
	}
	fmt.Printf("%s ", code)
	// Assign the generated code to the author
	author.Code = code

	// Tạo tác giả mới
	if err := h.DB.Create(&author).Error; err != nil {
		// Log the detailed error from database creation
		fmt.Printf("Error creating author in DB: %s\n", err.Error())
		// Return the detailed error to the client
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create author in DB: %s", err.Error())})
		return
	}

	// Successfully created the author, return it in the response
	c.JSON(http.StatusCreated, author)
}

// Hàm lấy thông tin tác giả theo ID
func (h *AuthorHandler) GetAuthorByID(c *gin.Context) {
	id := c.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author models.Author
	if err := h.DB.Preload("Books").First(&author, authorID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve author"})
		}
		return
	}

	c.JSON(http.StatusOK, author)
}

// Hàm cập nhật thông tin tác giả
func (h *AuthorHandler) UpdateAuthor(c *gin.Context) {
	id := c.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author models.Author
	if err := h.DB.First(&author, authorID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve author"})
		}
		return
	}

	// Bind dữ liệu mới từ request body
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra nếu name bị bỏ trống
	if strings.TrimSpace(author.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author name is required"})
		return
	}

	// Cập nhật thông tin tác giả
	if err := h.DB.Save(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// Hàm cập nhật thông tin tác giả (PATCH) - Chỉ cập nhật những trường được truyền
func (h *AuthorHandler) PatchAuthor(c *gin.Context) {
	id := c.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author models.Author
	if err := h.DB.First(&author, authorID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve author"})
		}
		return
	}

	// Tạo một struct tạm để chỉ bind các trường được truyền trong request
	var updatedAuthor struct {
		Name string `json:"name"`
	}

	// Bind dữ liệu JSON từ request body
	if err := c.ShouldBindJSON(&updatedAuthor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	// Kiểm tra nếu tên tác giả bị bỏ trống
	if strings.TrimSpace(updatedAuthor.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author name is required"})
		return
	}

	// Cập nhật các trường được truyền trong request
	if updatedAuthor.Name != "" {
		author.Name = updatedAuthor.Name
	}

	// Lưu thông tin đã cập nhật
	if err := h.DB.Save(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.JSON(http.StatusOK, author)
}

// Hàm xóa tác giả
func (h *AuthorHandler) DeleteAuthor(c *gin.Context) {
	id := c.Param("id")
	authorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid author ID"})
		return
	}

	var author models.Author
	if err := h.DB.First(&author, authorID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve author"})
		}
		return
	}

	// Xóa tác giả
	if err := h.DB.Delete(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	c.JSON(http.StatusNoContent, nil) // Trả về 204 No Content sau khi xóa thành công
}
