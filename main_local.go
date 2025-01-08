package main

import (
	"shop-account/handlers"
	"shop-account/routes"
	"shop-account/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func init() {
	// Kết nối tới database
	var err error
	DB, err = gorm.Open("postgres", "host=localhost user=postgres dbname=book_store password=1234 sslmode=disable")
	// DB, err = gorm.Open("postgres", "host=hotel-dqh20317-8f11.b.aivencloud.com user=avnadmin dbname=defaultdb password=AVNS_vqxx-jTp62srIABmstw port=25696 sslmode=disable")

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		os.Exit(1)
	}

	// Tự động migrate các model vào database
	DB.AutoMigrate(&models.Author{}, &models.Book{})
}

func main() {
	// Khởi tạo Gin router
	r := gin.Default()

	// Khởi tạo các handler
	authorHandler := &handlers.AuthorHandler{DB: DB}
	bookHandler := &handlers.BookHandler{DB: DB}
	authHandler := &handlers.AuthHandler{} // Không cần DB trong AuthHandler

	// Đăng ký tất cả các route
	routes.SetupRoutes(r, authorHandler, bookHandler, authHandler)

	// Khởi động server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
