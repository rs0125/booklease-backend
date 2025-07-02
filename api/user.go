package api

import (
	"bookapi/models"
	"bookapi/services"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func CreateOrFetchUser(c *gin.Context) {
	uid := c.GetString("uid")
	name := c.GetString("name") // Optional: you can also set this in middleware if needed

	if uid == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user data from token"})
		return
	}

	regNumPattern := regexp.MustCompile(`^\d{2}[A-Z]{3}\d{4}$`)
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

func GetUserProfile(c *gin.Context) {
	uid := c.GetString("uid") // assuming middleware has already set UID in context

	var user models.User
	result := services.DB.Select("username", "phone_number", "registration_no", "is_admin").
		Where("uid = ?", uid).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// return only selected fields
	c.JSON(http.StatusOK, gin.H{
		"username":        user.Username,
		"phone_number":    user.PhoneNumber,
		"registration_no": user.RegistrationNo,
		"is_admin":        user.IsAdmin,
		"email":           c.GetString("email"),
	})
}

func UpdatePhoneNumber(c *gin.Context) {
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UID not found in context"})
		return
	}

	var req PhoneUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number in request"})
		return
	}

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
