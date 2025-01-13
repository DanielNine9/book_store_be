package handlers

import (
	"fmt"
	"shop-account/models"
	"shop-account/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"errors" 
)

// AuthHandler struct with DB field
type AuthHandler struct {
	DB *gorm.DB
}
// Register handler for user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User

	// Bind JSON input to user model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Check if the username already exists
	var existingUser models.User
	err := h.DB.Where("username = ?", user.Username).First(&existingUser).Error

	if err == nil {
		// Username already taken
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	} 

	// Check if the error is due to a missing record (username not found)
	if err != gorm.ErrRecordNotFound {
		// Handle other database errors (e.g., connection issues)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while checking username"})
		return
	}

	// Hash the password using bcrypt
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Set default role for the new user
	user.Role = "guest" // Default role, can be changed based on your application's logic

	// Save the user to the database
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login handler for user login
func (h *AuthHandler) Login(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var existingUser models.User
	if err := h.DB.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !existingUser.Active {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User account is not active"})
		return
	}

	if !utils.ComparePasswords(existingUser.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(existingUser.ID, existingUser.Username, existingUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"id":       existingUser.ID,
		"username": existingUser.Username,
		// "email":    existingUser.Email,  // Include the email if necessary
		"role":     existingUser.Role,
		"active":   existingUser.Active, // Include active status if needed
	})
}


func generateToken(user_id uint, username, role string) (string, error) {
	secretKey := []byte("your_secret_key")
	fmt.Printf("user_id in generateToken %d", user_id)
	// the JWT claims require the value to be of type interface{}
	//  (which can hold various types), the correct approach would be to convert the uint to float64 
	// before setting it in the JWT claims. This is because JWT claims often store numbers as float64.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(user_id),
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(secretKey)
}

func (h *AuthHandler) ValidateToken(c *gin.Context) (*models.User, error) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("Missing authorization token")
	}

	// Extract token from the "Bearer <token>" format
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return nil, errors.New("Invalid authorization format")
	}

	// Parse the token
	secretKey := []byte("your_secret_key")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid token claims")
	}

	// Retrieve username and role from token claims
	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("Invalid username in token")
	}

	// Retrieve the user from the database
	var user models.User
	if err := h.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

// UpdateRole handler for admins to change the role of a user
func (h *AuthHandler) UpdateRole(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	// Bind JSON input to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	fmt.Printf("Username %s", request.Username)
	// Validate the JWT token and extract the user
	currentUser, err := h.ValidateToken(c)
	if err != nil || currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can change roles"})
		return
	}

	// Retrieve the user to update from the database
	
	var user models.User
	if err := h.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update the role
	user.Role = request.Role
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// ActivateUser handler for activating a user's account
func (h *AuthHandler) ActivateUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	// Bind the incoming request JSON to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate the JWT token and extract the user (assuming admin role is required)
	currentUser, err := h.ValidateToken(c)
	if err != nil || currentUser.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can activate users"})
		return
	}

	// Retrieve the user from the database
	var user models.User
	if err := h.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if the user is already active
	if user.Active {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already active"})
		return
	}

	// Update the user's Active status to true
	user.Active = true
	if err := h.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "User activated successfully"})
}
