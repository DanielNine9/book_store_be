package models

import "github.com/jinzhu/gorm"

type Book struct {
    gorm.Model
    Title          string  `json:"title"`
    Description    string  `json:"description"`
    Price          float64 `json:"price" gorm:"default:0"` 
    AuthorID       uint    `json:"author_id"`
    Author         Author  `json:"author"`
    Active         bool    `json:"active" gorm:"default:true"`
    QuantityInStock uint   `json:"quantity_in_stock" gorm:"default:10"` 
    QuantitySold   uint    `json:"quantity_sold" gorm:"default:0"` 
    Category         Category `json:"category"`   
    CategoryID       uint    `json:"category_id"` 
}
