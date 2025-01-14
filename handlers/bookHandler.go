package handlers

import (
	"shop-account/models"
	"shop-account/utils"
	"shop-account/dtos"
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

// func (h *BookHandler) GetBooks(c *gin.Context) {
// 	// Initialize an empty slice to hold the books
// 	var books []models.Book

// 	// Preload related data (Author and Category) separately
// 	query := h.DB.Preload("Author").Preload("Categories")

// 	// Call PaginateAndSearch utility to fetch paginated data with dynamic search (if any)
// 	totalItems, page, totalPages, err := utils.PaginateAndSearch(c, query, &models.Book{}, &books, nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books", "details": err.Error()})
// 		return
// 	}

// 	// Returning the paginated books as a response
// 	c.JSON(http.StatusOK, gin.H{
// 		"current_page":   page,
// 		"total_pages":    totalPages,
// 		"total_items":    totalItems,
// 		"items_per_page": c.DefaultQuery("limit", "10"),
// 		"books":          books,
// 	})
// }

func (h *BookHandler) GetBooks(c *gin.Context) {
	var books []models.Book

	// Query to preload related Author and Categories
	query := h.DB.Preload("Author").Preload("Categories")

	// Paginate and search
	totalItems, page, totalPages, err := utils.PaginateAndSearch(c, query, &models.Book{}, &books, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books", "details": err.Error()})
		return
	}

	// Prepare the response struct for books
	var bookResponses []dtos.BookResponse

	// Get user ID from context
	userID, exists := c.Get("id")
	if !exists {
		// If userID doesn't exist, we proceed with no favorites info
		for _, book := range books {
			bookResponses = append(bookResponses, dtos.BookResponse{
				ID:          book.ID,
				Title:       book.Title,
				Description: book.Description,
				Price:       book.Price,
				Author:      book.Author,
				Categories:  book.Categories,
				IsFavorite:  false,
				IdFavorite:  0,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"current_page":   page,
			"total_pages":    totalPages,
			"total_items":    totalItems,
			"items_per_page": c.DefaultQuery("limit", "10"),
			"books":          bookResponses,
			"userID":         userID,
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Loop through the books and build the response
	for _, book := range books {
		var favoriteBook models.FavoriteBook
		isFavorite := false
		var favoriteID *uint // Pointer to store favorite book ID

		// Check if the book is in the user's favorites
		if err := h.DB.Where("user_id = ? AND book_id = ?", userIDUint, book.ID).First(&favoriteBook).Error; err == nil {
			isFavorite = true
			favoriteID = &favoriteBook.ID // Assign the favorite book's ID if found
		}

		// Safely dereference favoriteID if it's not nil
		var favoriteIDValue uint
		if favoriteID != nil {
			favoriteIDValue = *favoriteID // Dereference to get the actual ID
		}

		// Build the book response
		bookResponse := dtos.BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Description: book.Description,
			Price:       book.Price,
			Author:      book.Author,
			Categories:  book.Categories,
			IsFavorite:  isFavorite,
			IdFavorite:  favoriteIDValue, // Set the ID of the favorite (or 0 if not found)
		}

		bookResponses = append(bookResponses, bookResponse)
	}

	// Return the response with books and pagination info
	c.JSON(http.StatusOK, gin.H{
		"current_page":   page,
		"total_pages":    totalPages,
		"total_items":    totalItems,
		"items_per_page": c.DefaultQuery("limit", "10"),
		"books":          bookResponses,
		"userID":         userID,
	})
}




func (h *BookHandler) CreateBook(c *gin.Context) {
    var requestData struct {
        Title       string   `json:"title"`
		Price       uint   `json:"price"`
		QuantityInStock       uint   `json:"quantity"`
        Description string   `json:"description"`
        AuthorID    uint     `json:"author_id"`
        CategoryIDs []uint   `json:"categories"`
    }

    if err := c.ShouldBindJSON(&requestData); err != nil {
        fmt.Printf("Error binding JSON: %s\n", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error binding JSON: %s", err.Error())})
        return
    }

    if strings.TrimSpace(requestData.Title) == "" {
        errMsg := "Book title is required"
        fmt.Printf("%s\n", errMsg)
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
        return
    }

    var author models.Author
    if err := h.DB.First(&author, requestData.AuthorID).Error; err != nil {
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

    book := models.Book{
        Title:       requestData.Title,
        Description: requestData.Description,
        AuthorID:    requestData.AuthorID,
        Price:    float64(requestData.Price),
        QuantityInStock:    requestData.QuantityInStock,
        Code:        code,
        Active:      true,
    }

    if len(requestData.CategoryIDs) > 0 {
        var categories []models.Category
        if err := h.DB.Find(&categories, "id IN (?)", requestData.CategoryIDs).Error; err != nil {
            fmt.Printf("Error finding categories: %s\n", err.Error())
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find categories"})
            return
        }

        book.Categories = categories
    } else {
        errMsg := "At least one category is required"
        fmt.Printf("%s\n", errMsg)
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
        return
    }

    if err := h.DB.Preload("Author").Preload("Categories").Create(&book).Error; err != nil {
        fmt.Printf("Error creating book in DB: %s\n", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create book in DB: %s", err.Error())})
        return
    }

    c.JSON(http.StatusCreated, book)
}
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id := c.Param("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := h.DB.Preload("Author").Preload("Categories").First(&book, bookID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	// Retrieve userID from context
	userID, exists := c.Get("id")
	var isFavorite bool
	var favoriteID uint // This will hold the favorite book ID

	if exists {
		userIDUint, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
			return
		}

		// Check if the book is in the user's favorites
		var favoriteBook models.FavoriteBook
		if err := h.DB.Where("user_id = ? AND book_id = ?", userIDUint, book.ID).First(&favoriteBook).Error; err == nil {
			isFavorite = true
			favoriteID = favoriteBook.ID // Use the ID directly, no need for dereferencing
		}
	}

	// Prepare the response struct
	bookResponse := dtos.BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Price:       book.Price,
		Author:      book.Author,
		Categories:  book.Categories,
		IsFavorite:  isFavorite,
		IdFavorite:  favoriteID, // The ID of the favorite or 0 if not a favorite
	}

	// Return the response
	c.JSON(http.StatusOK, bookResponse)
}


func (h *BookHandler) UpdateBook(c *gin.Context) {
    var requestData struct {
        Title            string   `json:"title"`
        Price            uint     `json:"price"`
        QuantityInStock  uint     `json:"quantity"`
        Description      string   `json:"description"`
        AuthorID         uint     `json:"author_id"`
        CategoryIDs      []uint   `json:"categories"` 
    }

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

    if err := c.ShouldBindJSON(&requestData); err != nil {
        fmt.Printf("Error binding JSON: %s\n", err.Error())
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error binding JSON: %s", err.Error())})
        return
    }

    if strings.TrimSpace(requestData.Title) == "" {
        errMsg := "Book title is required"
        fmt.Printf("%s\n", errMsg)
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
        return
    }

    var author models.Author
    if err := h.DB.First(&author, requestData.AuthorID).Error; err != nil {
        if gorm.IsRecordNotFoundError(err) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate author"})
        }
        return
    }

    if err := h.DB.Where("book_id = ?", book.ID).Delete(&models.BookCategory{}).Error; err != nil {
        fmt.Printf("Error manually clearing categories: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to manually clear categories"})
        return
    }

    book.Title = requestData.Title
    book.Description = requestData.Description
    book.Price = float64(requestData.Price)
    book.QuantityInStock = requestData.QuantityInStock
    book.AuthorID = requestData.AuthorID

    if len(requestData.CategoryIDs) > 0 {
        var categories []models.Category
        if err := h.DB.Find(&categories, "id IN (?)", requestData.CategoryIDs).Error; err != nil {
            fmt.Printf("Error finding categories: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find categories"})
            return
        }
		// if err := h.DB.Model(&book).Association("Categories").Append(categories); err != nil {
        //     fmt.Printf("Error updating categories: %v\n", err)
        //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update categories"})
        //     return
        // }
        // Insert new relationships into the book_categories join table
        for _, category := range categories {
            bookCategory := models.BookCategory{
                BookID:     book.ID,
                CategoryID: category.ID,
            }
            if err := h.DB.Create(&bookCategory).Error; err != nil {
                fmt.Printf("Error creating book-category relationship: %v\n", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update categories"})
                return
            }
        }
    } else {
        errMsg := "At least one category is required"
        fmt.Printf("%s\n", errMsg)
        c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
        return
    }

    if err := h.DB.Preload("Author").Preload("Categories").Save(&book).Error; err != nil {
        fmt.Printf("Error updating book in DB: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
        return
    }

    c.JSON(http.StatusOK, book)
}




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

	if err := h.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusNoContent, nil) 
}

func (h *BookHandler) PatchBook(c *gin.Context) {
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

	var updatedBook models.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

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

func (h *BookHandler) Restore(c *gin.Context) {
    id := c.Param("id")
    var book models.Book

    if err := h.DB.Unscoped().Where("id = ?", id).First(&book).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    if book.DeletedAt == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Book is not deleted"})
        return
    }

    if err := h.DB.Model(&book).Update("deleted_at", nil).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore book"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book restored successfully"})
}


