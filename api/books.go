package api

import (
	"net/http"

	"bookapi/models"

	"github.com/gin-gonic/gin"
)

var books = []models.Book{
	{ID: "1", Title: "Go in Action", Author: "William Kennedy"},
	{ID: "2", Title: "The Go Programming Language", Author: "Alan Donovan"},
}

func GetBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	for _, b := range books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func CreateBook(c *gin.Context) {
	var newBook models.Book
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}
