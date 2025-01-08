package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6"`
	Role string `json:"role"`
}
