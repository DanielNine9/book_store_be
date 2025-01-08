package main

import (
	"log"
	"os"
	"shop-account/models"
	"shop-account/handlers"
	"shop-account/routes"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB *gorm.DB

func init() {
	var err error

	// Connection string to Aiven PostgreSQL database with SSL
	// serviceURI := "postgres://avnadmin:AVNS_vqxx-jTp62srIABmstw@hotel-dqh20317-8f11.b.aivencloud.com:25696/defaultdb?sslmode=require"
	serviceURI := "postgres://postgres:1234@localhost:5432/book_store?sslmode=disable"

	// Open connection to PostgreSQL using GORM
	DB, err = gorm.Open("postgres", serviceURI)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		os.Exit(1)
	}

	// Test the connection (this ensures the database is accessible)
	if err := DB.DB().Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
		os.Exit(1)
	}

	if err := DB.AutoMigrate(&models.Author{}, &models.Book{}, &models.User{}).Error; err != nil {
		log.Fatal("Failed to migrate database:", err)
		os.Exit(1)
	}

	log.Println("Successfully connected to the database")
}

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize handlers
	authorHandler := &handlers.AuthorHandler{DB: DB}
	bookHandler := &handlers.BookHandler{DB: DB}
	authHandler := &handlers.AuthHandler{DB: DB} // No need for DB in AuthHandler
	userHandler := &handlers.UserHandler{DB: DB} 
	// Register routes
	routes.SetupRoutes(r,userHandler, authorHandler, bookHandler, authHandler, )

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
