package models

import "github.com/jinzhu/gorm"

type Book struct {
    gorm.Model
    Title       string `json:"title"`
    Description string `json:"description"`
    AuthorID    uint   `json:"author_id"`
    Author      Author `json:"author"`
    
	Active   bool   `json:"active" gorm:"default:false"` 
}
