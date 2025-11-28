package internal

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecom/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler holds dependencies (like the DB connection)
type AuthHandler struct {
	DB *gorm.DB
}

// RegisterRequest defines what the frontend must send us
// Binding tags (json, binding) ensure we get valid data automatically
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register is the actual function that handles the route
func (h *AuthHandler) Register(c *gin.Context) {
	var body RegisterRequest

	// 1. Validation: Parse JSON and check rules (email format, min length)
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hash the Password
	// Cost 14 is a good balance between security and speed
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Create User Model
	user := models.User{
		Email:    body.Email,
		Password: string(hash), // Store the HASH, not the password
	}

	// 4. Save to Database
	result := h.DB.Create(&user)

	// Check for duplicate email error (Unique Constraint)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// 5. Success Response
	// NEVER return the password (even the hash) back to the user
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"userId":  user.ID,
		"email":   user.Email,
	})
}

// LoginRequest defines what we need to log in
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body LoginRequest

	// 1. Validation
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Find User by Email
	var user models.User
	result := h.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 3. Verify Password
	// Compare the stored Hash with the plain text password sent by user
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 4. Generate JWT Token
	// In a real app, load this secret from os.Getenv("JWT_SECRET")
	secretKey := []byte("supersecretkey")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,                               // Subject (User ID)
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration (24 hours)
	})

	// Sign the token with our secret
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 5. Respond
	// We send the token back. The frontend will store this.
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
