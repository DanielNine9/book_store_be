package handlers

import (
	"shop-account/models"
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

// Hàm lấy danh sách sách
func (h *BookHandler) GetBooks(c *gin.Context) {
	var books []models.Book
	if err := h.DB.Preload("Author").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// Hàm tạo sách mới
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book
	// Bind dữ liệu JSON từ request body
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

	// Tạo sách mới
	if err := h.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
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
	if err := h.DB.Preload("Author").First(&book, bookID).Error; err != nil {
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
	if err := h.DB.First(&author, updatedBook.AuthorID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate author"})
		}
		return
	}

	if updatedBook.Title != "" {
		book.Title = updatedBook.Title
	}
	if updatedBook.AuthorID != 0 {
		book.AuthorID = updatedBook.AuthorID
	}

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

// Hàm lấy danh sách sách và tác giả đồng thời (sử dụng concurrency)
func (h *BookHandler) GetBooksConcurrently(c *gin.Context) {
	// Tạo hai channel để nhận kết quả từ các goroutines
	bookResultChan := make(chan BookResult)
	authorResultChan := make(chan AuthorResult)

	// Goroutine để lấy danh sách sách
	go func() {
		var books []models.Book
		err := h.DB.Preload("Author").Find(&books).Error
		bookResultChan <- BookResult{Books: books, Err: err} // Gửi kết quả vào channel bookResultChan
	}()

	// Goroutine để lấy danh sách tác giả
	go func() {
		var authors []models.Author
		err := h.DB.Find(&authors).Error
		authorResultChan <- AuthorResult{Authors: authors, Err: err} // Gửi kết quả vào channel authorResultChan
	}()

	// Chờ và nhận kết quả từ cả hai channel
	var books []models.Book
	var authors []models.Author

	// Nhận kết quả từ cả hai channel
	for i := 0; i < 2; i++ {
		select {
		case result := <-bookResultChan:
			if result.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve books: %v", result.Err)})
				return
			}
			books = result.Books // Lưu trữ kết quả sách
		case result := <-authorResultChan:
			if result.Err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve authors: %v", result.Err)})
				return
			}
			authors = result.Authors // Lưu trữ kết quả tác giả
		}
	}

	// Trả về kết quả sau khi đã lấy xong dữ liệu cả sách và tác giả
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
// Hàm import 1000 sách vào cơ sở dữ liệu
func (h *BookHandler) ImportBooks(c *gin.Context) {
	// Kiểm tra xem tác giả có tồn tại không
	var author models.Author
	if err := h.DB.First(&author, 1).Error; err != nil {
		// Nếu không có tác giả với ID = 1, trả về lỗi
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	// Tạo dữ liệu giả cho 1000 sách
	var books []models.Book
	for i := 1; i <= 10000; i++ {
		book := models.Book{
			Title:       fmt.Sprintf("Book Title %d", i),
			Description: fmt.Sprintf("Description for Book %d", i),
			AuthorID:    author.ID, // Đảm bảo sử dụng ID của tác giả tồn tại
		}
		books = append(books, book)
	}

	// Log mảng books để kiểm tra
	log.Printf("Books to import: %+v", books)

	// Kiểm tra mảng books có trống không
	if len(books) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No books to import"})
		return
	}

	
	// Thêm các sách vào cơ sở dữ liệu
		for _, book := range books {
		if err := h.DB.Create(&book).Error; err != nil {
			// Nếu gặp lỗi khi thêm một cuốn sách, log lỗi và tiếp tục với sách tiếp theo
			log.Printf("Error importing book: %+v. Error: %v", book, err)
			continue
		}
		// Log thành công cho mỗi cuốn sách
		log.Printf("Successfully imported book: %+v", book)
	}

	// Trả về kết quả sau khi import thành công
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully imported %d books", len(books)),
	})
}
