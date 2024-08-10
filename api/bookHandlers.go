package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/dataModels"
	"shop/db"
	"shop/logger"
	"strconv"
)

func addBook(c *gin.Context) {
	var newBook dataModels.Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if newBook.Title == "" || newBook.ISBN == "" || newBook.Language == "" || newBook.Year == 0 {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("missing required fields"))
		return
	}
	if err := db.AddBook(newBook); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book added successfully!"})
}

func getBooks(c *gin.Context) {
	quantityStr := c.Param("quantity")

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid quantity"))
		return
	}
	if quantity <= 0 {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("the quantity should not be negative"))
		return
	}
	books, err := db.GetBooks(quantity)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, books)
}

func getBookById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}
	book, err := db.GetBookById(id)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("book not found"))
		return
	}
	if book.ISBN == "" && book.Title == "" {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("book not found"))
		return
	}

	c.JSON(http.StatusOK, book)
}

func searchByTitle(c *gin.Context) {
	title := c.Param("title")
	books, err := db.SearchByTitle(title)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
	}

	if len(books) == 0 {
		errorResponse(c, http.StatusNotFound, fmt.Errorf("record not found"))
		return
	}
	c.JSON(http.StatusOK, books)
}

func errorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error:": err.Error()})
	logger.Error(err)
}
