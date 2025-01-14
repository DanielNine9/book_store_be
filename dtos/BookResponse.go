
package dtos


import (
    "shop-account/models"
)

type BookResponse struct {
    ID            uint       `json:"id"`
    Title         string     `json:"title"`
    Description   string     `json:"description"`
    Price         float64    `json:"price"`
    Author        models.Author     `json:"author"`
    BookDetail        models.BookDetail     `json:"book_detail"`
    Images        []models.BookImage     `json:"images"`
    Categories    []models.Category `json:"categories"`
    IsFavorite    bool       `json:"is_favorite"` 
    IdFavorite    uint   `json:"id_favorite"`
    Active          bool       `json:"active"`
    QuantityInStock uint       `json:"quantity_in_stock""`
    QuantitySold    uint       `json:"quantity_sold"`
}