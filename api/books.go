package api

import (
	"net/http"
	"strconv"

	"bookapi/models"
	"bookapi/services"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	var books []models.Book
	if err := services.DB.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var book models.Book
	if err := services.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.DB.Create(&newBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newBook)
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := services.DB.Delete(&models.Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func AddToWishlist(c *gin.Context) {

	bookIDStr := c.Param("id")
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var requestBody struct {
		UserID uint `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil || requestBody.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid user_id"})
		return
	}

	// Check if book exists
	var book models.Book
	if err := services.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var existing models.Wishlist
	if err := services.DB.
		Where("user_id = ? AND book_id = ?", requestBody.UserID, bookID).
		First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Book already in wishlist"})
		return
	}

	wish := models.Wishlist{
		UserID: requestBody.UserID,
		BookID: uint(bookID),
	}
	if err := services.DB.Create(&wish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book added to wishlist"})
}

func Wishlist(c *gin.Context) {
	var wishlistBooks []models.Book
	if err := services.DB.Find(&wishlistBooks).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
		return
	}
	c.IndentedJSON(http.StatusOK, wishlistBooks)
}
