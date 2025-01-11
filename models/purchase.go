package models

import "github.com/jinzhu/gorm"

type Purchase struct {
    gorm.Model
    UserID      uint    `json:"user_id"`
    BookID      uint    `json:"book_id"`
    Quantity    uint    `json:"quantity" gorm:"default:1"`  
    User        User    `json:"user"`
    Book        Book    `json:"book"`
    TransactionID uint   `json:"transaction_id"` 
    BookPrice   float64 `json:"book_price" gorm:"default:1`     
    Transaction Transaction `json:"transaction"`    
     Code        string `json:"code"`
}
