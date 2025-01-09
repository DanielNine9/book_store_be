package dtos

import (
    "shop-account/models"
)

type PurchaseResponse struct {
    ID        uint      `json:"id"`
    CreatedAt string    `json:"created_at"`
    UpdatedAt string    `json:"updated_at"`
    DeletedAt *string `json:"deleted_at"`
    UserID    uint      `json:"user_id"`
    Book    models.Book      `json:"book"`
    Quantity  uint      `json:"quantity"`
}