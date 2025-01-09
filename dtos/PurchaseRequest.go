package dtos

type PurchaseRequest struct {
    Quantity uint `json:"quantity" binding:"required"`
}