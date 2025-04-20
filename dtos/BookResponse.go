
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
    Categories    []models.Category `json:"categories"`
    IsFavorite    bool       `json:"is_favorite"` 
    IdFavorite    uint   `json:"id_favorite"`
}