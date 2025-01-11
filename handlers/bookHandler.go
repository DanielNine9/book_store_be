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
	"log"
)

// Handler struct để lưu trữ db instance
type BookHandler struct {
	DB *gorm.DB
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	// Initialize an empty slice to hold the books
	var books []models.Book

	// Preload related data (Author and Category) separately
	query := h.DB.Preload("Author").Preload("Category")

	// Call PaginateAndSearch utility to fetch paginated data with dynamic search (if any)
	totalItems, page, totalPages, err := utils.PaginateAndSearch(c, query, &models.Book{}, &books, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books", "details": err.Error()})
		return
	}

	// Returning the paginated books as a response
	c.JSON(http.StatusOK, gin.H{
		"current_page":   page,
		"total_pages":    totalPages,
		"total_items":    totalItems,
		"items_per_page": c.DefaultQuery("limit", "10"),
		"books":          books,
	})
}


// Hàm tạo sách mới
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		fmt.Printf("Error binding JSON: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error binding JSON: %s", err.Error())})
		return
	}

	if strings.TrimSpace(book.Title) == "" {
		errMsg := "Book title is required"
		fmt.Printf("%s\n", errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	var author models.Author
	if err := h.DB.First(&author, book.AuthorID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate author"})
		}
		return
	}

	code, err := utils.GenerateCode(h.DB, &models.Book{})
	if err != nil {
		fmt.Printf("Error generating book code: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate book code: %s", err.Error())})
		return
	}

	book.Code = code
	book.Active = true
	if err := h.DB.Preload("Author").Preload("Category").Create(&book).Error; err != nil {
		fmt.Printf("Error creating book in DB: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create book in DB: %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, book)
}


// Hàm lấy thông tin sách theo ID
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := h.DB.Preload("Author").Preload("Category").First(&book, bookID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// Hàm cập nhật thông tin sách
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := h.DB.First(&book, bookID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	// Bind dữ liệu mới từ request body
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Kiểm tra nếu tên sách bị bỏ trống
	if strings.TrimSpace(book.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book title is required"})
		return
	}

	// Kiểm tra nếu tác giả không hợp lệ (ví dụ: không tồn tại)
	var author models.Author
	if err := h.DB.First(&author, book.AuthorID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate author"})
		}
		return
	}

	// Cập nhật thông tin sách
	if err := h.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Hàm xóa sách
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := h.DB.First(&book, bookID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	// Xóa sách
	if err := h.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusNoContent, nil) // Trả về 204 No Content sau khi xóa thành công
}

func (h *BookHandler) PatchBook(c *gin.Context) {
	// Lấy ID từ URL params
	id := c.Param("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Tìm sách theo ID
	var book models.Book
	if err := h.DB.First(&book, bookID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	// Bind dữ liệu PATCH từ request body vào book
	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Kiểm tra nếu tác giả không hợp lệ (ví dụ: không tồn tại)
	var author models.Author
	if updatedBook.AuthorID != 0 {
		if err := h.DB.First(&author, updatedBook.AuthorID).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate author"})
			}
			return
		}
	}
	

	if updatedBook.Title != "" {
		book.Title = updatedBook.Title
	}
	if updatedBook.AuthorID != 0 {
		book.AuthorID = updatedBook.AuthorID
	}
	if updatedBook.Price != 0 {
		book.Price = updatedBook.Price
	}
	if updatedBook.QuantityInStock != 0 {
		book.QuantityInStock = updatedBook.QuantityInStock
	}
		
	book.Active = updatedBook.Active

	if err := h.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}
type BookResult struct {
	Books []models.Book
	Err   error
}

// AuthorResult để chứa dữ liệu kết quả và lỗi trả về từ channel
type AuthorResult struct {
	Authors []models.Author
	Err     error
}

func (h *BookHandler) GetBooksConcurrently(c *gin.Context) {
	bookResultChan := make(chan BookResult)
	authorResultChan := make(chan AuthorResult)

	go func() {
		var books []models.Book
		err := h.DB.Preload("Author").Find(&books).Error
		bookResultChan <- BookResult{Books: books, Err: err} 
	}()

	go func() {
		var authors []models.Author
		err := h.DB.Find(&authors).Error
		authorResultChan <- AuthorResult{Authors: authors, Err: err} 
	}()

	var books []models.Book
	var authors []models.Author

	for i := 0; i < 2; i++ {
		select {
		case result := <-bookResultChan:
			if result.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve books: %v", result.Err)})
				return
			}
			books = result.Books 
		case result := <-authorResultChan:
			if result.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve authors: %v", result.Err)})
				return
			}
			authors = result.Authors 
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"books":   books,
		"authors": authors,
	})
}


// Hàm lấy danh sách sách
func (h *BookHandler) GetBooksNotConcurrently(c *gin.Context) {
	var books []models.Book
	if err := h.DB.Preload("Author").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	var authors []models.Author
	h.DB.Preload("Books").Find(&authors)
	
	c.JSON(http.StatusOK, gin.H{"books": books, "authors": authors})
}


func (h *BookHandler) ImportBooks(c *gin.Context) {
	var author models.Author
	if err := h.DB.First(&author, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	var books []models.Book
	for i := 1; i <= 100; i++ {
		book := models.Book{
			Title:       fmt.Sprintf("Book Title %d", i),
			Description: fmt.Sprintf("Description for Book %d", i),
			AuthorID:    author.ID, 
		}
		books = append(books, book)
	}

	log.Printf("Books to import: %+v", books)

	if len(books) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No books to import"})
		return
	}

	
		for _, book := range books {
		if err := h.DB.Create(&book).Error; err != nil {
			log.Printf("Error importing book: %+v. Error: %v", book, err)
			continue
		}
		log.Printf("Successfully imported book: %+v", book)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully imported %d books", len(books)),
	})
}

// Restore function to restore a soft-deleted record
func (h *BookHandler) Restore(c *gin.Context) {
    id := c.Param("id")
    var book models.Book

    // Try to find the deleted book by its ID (including soft-deleted ones)
    if err := h.DB.Unscoped().Where("id = ?", id).First(&book).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    // Check if the book is already active (not deleted)
    if book.DeletedAt == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Book is not deleted"})
        return
    }

    // Restore the book by setting DeletedAt to nil
    if err := h.DB.Model(&book).Update("deleted_at", nil).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore book"})
        return
    }

    // Return success message
    c.JSON(http.StatusOK, gin.H{"message": "Book restored successfully"})
}


