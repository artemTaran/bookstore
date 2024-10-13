package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"shop/dataModels"
	"shop/db"
	"shop/logger"
	"strconv"
)

func getBooks(c *gin.Context) {
	quantityParam := c.Query("quantity")
	titleParam := c.Query("title")

	if titleParam != "" {
		getByTitle(c, titleParam)
		return
	}

	getListBooks(c, quantityParam)
}

func addBook(c *gin.Context) {
	var newBook dataModels.Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		errorResponse(c, http.StatusBadRequest, errors.New("there is something wrong with the data you sent"))
		logger.Error(err)
		return
	}

	if err := newBook.Validate(); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := db.AddBook(newBook); err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "book added successfully!"})
}

func getListBooks(c *gin.Context, quantityStr string) {
	if quantityStr == "" {
		quantityStr = "10"
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("the quantity must be number"))
		return
	}
	if quantity <= 0 {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("the quantity should not be negative"))
		return
	}
	books, err := db.GetBooks(quantity)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, errors.New("error getting list of books from database, try again later"))
		logger.Error(err)
		return
	}

	booksResponse := createResponseBooks(books)

	c.JSON(http.StatusOK, booksResponse)
}

func getBookById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("the id param must be a number"))
		return
	}
	book, err := db.GetBookById(id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, fmt.Errorf("book not found"))
		return
	}
	if book.ISBN == "" && book.Title == "" {
		errorResponse(c, http.StatusNotFound, fmt.Errorf("book not found"))
		return
	}

	bookResponse := createResponseBooks([]dataModels.Book{book})

	c.JSON(http.StatusOK, bookResponse)
}

func getByTitle(c *gin.Context, title string) {
	books, err := db.SearchByTitle(title)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
	}

	if len(books) == 0 {
		errorResponse(c, http.StatusNotFound, fmt.Errorf("record not found"))
		return
	}

	booksResponse := createResponseBooks(books)

	c.JSON(http.StatusOK, booksResponse)
}

func errorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error:": err.Error()})
	logger.Error(err)
}

func createResponseBooks(books []dataModels.Book) []dataModels.BookResponse {
	bookResponses := make([]dataModels.BookResponse, len(books))

	for i, book := range books {
		authors := make([]dataModels.AuthorResponse, len(book.Authors))
		for j, author := range book.Authors {
			authors[j] = dataModels.AuthorResponse{
				ID:        author.ID,
				FirstName: author.FirstName,
				LastName:  author.LastName,
			}
		}

		bookResponses[i] = dataModels.BookResponse{
			ID:          book.ID,
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     authors,
			Description: book.Description,
			Language:    book.Language,
			Year:        book.Year,
		}
	}
	return bookResponses
}
