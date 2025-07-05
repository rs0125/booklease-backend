package api

import (
	"bookapi/models"
	"bookapi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// import (
// 	"net/http"
// 	"strconv"

// 	"bookapi/models"
// 	"bookapi/services"

// 	"github.com/gin-gonic/gin"
// )

func PostRental(c *gin.Context) {

	uid := c.GetString("uid")
	if uid == " " {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var newRental models.Rental
	if err := c.ShouldBindJSON(&newRental); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newRental.BookID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing book_id"})
		return
	}

	fmt.Printf("📦 Incoming rental: %+v\n", newRental)

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var book models.Book
	if err := services.DB.First(&book, *newRental.BookID).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	newRental.UserID = user.ID
	newRental.OwnerID = &book.UploadedBy
	newRental.IsReturned = false

	fmt.Printf("Final rental to save: %+v\n", newRental)

	if err := services.DB.Create(&newRental).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	services.CreateNotification(*&book.UploadedBy, "rental_request", user.Username+" wants to rent \""+book.Title+"\"")

	c.IndentedJSON(http.StatusCreated, newRental)

}

func GetRentals(c *gin.Context) {

	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var rentals []models.Rental

	if err := services.DB.Find(&rentals).Order("ID DESC").Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, rentals)
}

func BorrowedMaterials(c *gin.Context) {
	uid := c.GetString("uid")

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var rentals []models.Rental
	if err := services.DB.Where("user_id = ?", user.ID).Order("ID DESC").Preload("Book").Find(&rentals).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, rentals)
}

func LentMaterials(c *gin.Context) {
	uid := c.GetString("uid")

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var rentals []models.Rental
	if err := services.DB.
		Where("owner_id = ?", user.ID). //add this for excluding self rentals
		Preload("Book").
		Preload("Book.Uploader").
		Preload("User"). // renter info

		Find(&rentals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lent materials"})
		return
	}

	c.JSON(http.StatusOK, rentals)
}

func DecideRental(c *gin.Context) {
	uid := c.GetString("uid")
	rentalID := c.Param("id")

	var user models.User
	if err := services.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var rental models.Rental
	if err := services.DB.Preload("Book").First(&rental, rentalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rental not found"})
		return
	}

	// Check ownership
	if rental.OwnerID == nil || *rental.OwnerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to decide this rental"})
		return
	}

	var body struct {
		Accept bool `json:"accept"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rental.Status = &body.Accept
	if err := services.DB.Save(&rental).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rental status"})
		return
	}

	statusText := "rejected"
	if body.Accept {
		statusText = "accepted"
	}

	services.CreateNotification(rental.UserID, "rental_status", "Your request to rent \""+rental.Book.Title+"\" was "+statusText)
	c.JSON(http.StatusOK, gin.H{"message": "Rental request " + statusText})
}
