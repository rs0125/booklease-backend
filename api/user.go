package api

import (
	"bookapi/models"
	"bookapi/services"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateOrFetchUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed Authorization header"})
		return
	}
	idToken := strings.TrimPrefix(authHeader, "Bearer ")

	// Firebase Auth client
	authClient, err := services.App.Auth(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase Auth init failed"})
		return
	}

	// Verify token
	token, err := authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
		return
	}

	// Extract name and email
	name, okName := token.Claims["name"].(string)
	_, okEmail := token.Claims["email"].(string)
	if !okName || !okEmail {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	// Parse registration number
	parts := strings.Fields(name)
	if len(parts) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse registration number from name"})
		return
	}
	registrationNo := parts[2]

	// Check if user exists
	var user models.User
	result := services.DB.Where("registration_no = ?", registrationNo).First(&user)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User already exists", "user": user})
		return
	}

	// Create new user
	newUser := models.User{
		Username:       name,
		RegistrationNo: registrationNo,
		IsAdmin:        false,
	}
	if err := services.DB.Create(&newUser).Error; err != nil {
		log.Println("❌ DB Create error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created", "user": newUser})
}
