package models

import "github.com/jinzhu/gorm"

type Author struct {
    gorm.Model
    Name  string `json:"name"`
    Bio   string `json:"bio"`
    Books []Book `json:"books"`
	Active   bool   `json:"active" gorm:"default:true"` 

}
