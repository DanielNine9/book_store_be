package main

import (
	"log"
	"os"
	"shop-account/models"
	"shop-account/handlers"
	"shop-account/handlers/admin"
	"shop-account/routes"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

var DB *gorm.DB

func init() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	serviceURI := os.Getenv("DATABASE_URL")
	if serviceURI == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	DB, err = gorm.Open("postgres", serviceURI)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		os.Exit(1)
	}

	if err := DB.DB().Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
		os.Exit(1)
	}

	if err := DB.AutoMigrate(&models.FavoriteBook{},&models.BookImage{},&models.BookCategory{}, &models.Category{}, &models.Author{}, &models.Book{}, &models.User{}, &models.Purchase{}, &models.Transaction{}).Error; err != nil {
		log.Fatal("Failed to migrate database:", err)
		os.Exit(1)
	}

	log.Println("Successfully connected to the database")
}

func main() {
	r := gin.Default()

	// Enable CORS with all origins, methods, and headers
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow all necessary HTTP methods
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Origin", "Accept"}, // Allow all relevant headers
		AllowCredentials: true,           // Allow credentials (cookies, authorization)
		MaxAge:           12 * 3600,     // Cache preflight response for 12 hours
	}))

	// Log the headers of incoming requests to confirm CORS headers
	r.Use(func(c *gin.Context) {
		log.Println("Request Headers:", c.Request.Header)
		c.Next()
	})

	// Initialize handlers
	authorHandler := &handlers.AuthorHandler{DB: DB}
	bookHandler := &handlers.BookHandler{DB: DB}
	authHandler := &handlers.AuthHandler{DB: DB}
	userHandler := &handlers.UserHandler{DB: DB}
	purchaseHandler := &handlers.PurchaseHandler{DB: DB}
	transactionHandler := &handlers.TransactionHandler{DB: DB}
	transactionAdminHandler := &admin.AdminTransactionHandler{DB: DB}
	categoryHandler := &handlers.CategoryHandler{DB: DB}
	favoriteHandler := &handlers.FavoriteBookHandler{DB: DB}

	// Set up routes
	routes.SetupRoutes(r, favoriteHandler, categoryHandler, transactionAdminHandler, transactionHandler, purchaseHandler, userHandler, authorHandler, bookHandler, authHandler)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
