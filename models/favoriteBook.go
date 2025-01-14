package models

import "github.com/jinzhu/gorm"

type FavoriteBook struct {
    gorm.Model
    UserID uint   `json:"user_id"`
    User   User   `json:"user"`
    BookID uint   `json:"book_id"`
    Book   Book   `json:"book"`
}
