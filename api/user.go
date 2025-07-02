package api

import (
	"bookapi/models"
	"bookapi/services"
	"context"
	"log"
	"net/http"
	"regexp"
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

	authClient, err := services.App.Auth(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase Auth init failed"})
		return
	}

	token, err := authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
		return
	}

	uid, okUID := token.Claims["user_id"].(string)
	name, okName := token.Claims["name"].(string)
	if !okUID || !okName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	var regNumPattern = regexp.MustCompile(`^\d{2}[A-Z]{3}\d{4}$`)

	parts := strings.Fields(name)
	var registrationNo string
	for _, word := range parts {
		if regNumPattern.MatchString(strings.ToUpper(word)) {
			registrationNo = strings.ToUpper(word)
			break
		}
	}

	if registrationNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to extract registration number"})
		return
	}
	var user models.User
	result := services.DB.Where("uid = ?", uid).First(&user)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"message": "User already exists", "user": user})
		return
	}

	newUser := models.User{
		UID:            uid,
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

type PhoneUpdateRequest struct {
	PhoneNumber string `json:"phone_number"`
}

func UpdatePhoneNumber(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed Authorization header"})
		return
	}
	idToken := strings.TrimPrefix(authHeader, "Bearer ")

	authClient, err := services.App.Auth(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase Auth init failed"})
		return
	}

	token, err := authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
		return
	}

	uid, ok := token.Claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UID not found in token"})
		return
	}

	var req PhoneUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number in request"})
		return
	}

	// Validate phone number using regex
	re := regexp.MustCompile(`^[6-9]\d{9}$`)
	if !re.MatchString(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number format invalid"})
		return
	}

	var user models.User
	result := services.DB.Where("uid = ?", uid).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.PhoneNumber = req.PhoneNumber
	if err := services.DB.Save(&user).Error; err != nil {
		log.Println("❌ DB Update error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update phone number"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Phone number updated", "user": user})
}
