package models

import "github.com/jinzhu/gorm"

type Category struct {
    gorm.Model
    Name        string `json:"name"` 
    Description string `json:"description"`
    ImageURL    string `json:"image_url"`
     Code        string `json:"code"`
}
