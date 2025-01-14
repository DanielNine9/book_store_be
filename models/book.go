package models

import "github.com/jinzhu/gorm"

type Book struct {
    gorm.Model
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    Price           float64   `json:"price" gorm:"default:0"`
    AuthorID        uint      `json:"author_id"`
    Author          Author    `json:"author"`
    Active          bool      `json:"active" gorm:"default:true"`
    QuantityInStock uint      `json:"quantity_in_stock" gorm:"default:10"`
    QuantitySold    uint      `json:"quantity_sold" gorm:"default:0"`
    Categories      []Category `gorm:"many2many:book_categories;foreignkey:ID;association_foreignkey:ID" json:"categories"`
    Code            string    `json:"code"`

    Publisher       string    `json:"publisher"`
    PublicationYear int       `json:"publication_year"`
    Weight          float64   `json:"weight"`  
    Dimensions      string    `json:"dimensions"`  
    Pages           int       `json:"pages"`  
    BindingType     string    `json:"binding_type"`  

    Images          []BookImage   `gorm:"foreignkey:ID" json:"images"`
}

type BookImage struct {
    gorm.Model
    URL     string `json:"url"`
    BookID  uint   `json:"book_id"`
}
