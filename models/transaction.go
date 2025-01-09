// package models

// import (
//     "time"
//     "github.com/jinzhu/gorm"
// )

// type Transaction struct {
//     gorm.Model
//     UserID          uint        `json:"user_id"`
//     TotalAmount     float64     `json:"total_amount"`
//     Status          string      `json:"status"`
//     TransactionTime time.Time   `json:"transaction_time"`
//     Purchases       []Purchase  `json:"purchases"`
//     User            User        `json:"user"`
	
// }
// func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
//     if t.TransactionTime.IsZero() {
//         t.TransactionTime = time.Now() 
//     }
//     return nil
// }

package models

import (
    "time"
    "github.com/jinzhu/gorm"
)

type TransactionStatus string

const (
    Pending   TransactionStatus = "pending"
    Approved  TransactionStatus = "approved"
    Rejected  TransactionStatus = "rejected"
    Completed TransactionStatus = "completed"
)

type Transaction struct {
    gorm.Model
    UserID          uint               `json:"user_id"`
    TotalAmount     float64            `json:"total_amount"`
    Status          TransactionStatus  `json:"status"`
    TransactionTime time.Time          `json:"transaction_time"`
    Purchases       []Purchase         `json:"purchases"`
    User            User               `json:"user"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
    if t.TransactionTime.IsZero() {
        t.TransactionTime = time.Now()
    }
    return nil
}
