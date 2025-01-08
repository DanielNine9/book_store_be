package models

import "github.com/jinzhu/gorm"

// Purchase struct represents the purchase record between user and book.
type Purchase struct {
    gorm.Model
    UserID uint   `json:"user_id"`
    BookID uint   `json:"book_id"`
    User   User   `json:"user"`
    Book   Book   `json:"book"`
}
