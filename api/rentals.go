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

	fmt.Printf("📝 Final rental to save: %+v\n", newRental)

	if err := services.DB.Create(&newRental).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	if err := services.DB.Find(&rentals).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, rentals)
}
