package models

import "github.com/jinzhu/gorm"

type Book struct {
    gorm.Model
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Price       float64 `json:"price" gorm:"default:0"` // Default value set to 0
    AuthorID    uint    `json:"author_id"`
    Author      Author  `json:"author"`
    
    Active      bool    `json:"active" gorm:"default:false"`  // Active status for the book
}
