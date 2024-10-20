package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"shop/dataModels"
	"shop/db"
	"shop/logger"
	"strconv"
)

func getAuthors(c *gin.Context) {
	quantityParam := c.Query("quantity")
	fullNameParam := c.Query("fullName")

	if fullNameParam != "" {
		getByFullName(c, fullNameParam)
		return
	}

	getListAuthors(c, quantityParam)
}

func getListAuthors(c *gin.Context, quantityStr string) {
	if quantityStr == "" {
		quantityStr = "10"
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, errors.New("the quantity must be number"))
		return
	}
	if quantity <= 0 {
		errorResponse(c, http.StatusBadRequest, errors.New("the quantity should not be negative"))
		return
	}
	authors, err := db.GetAuthors(quantity)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, errors.New("internal server error, please try making a request later"))
		logger.Error(err)
		return
	}

	authorsResponse := createResponseAuthors(authors)

	c.JSON(http.StatusOK, authorsResponse)
}

func getAuthorById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, errors.New("the id param must be a number"))
		return
	}
	author, err := db.GetAuthorById(id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, errors.New("author not found"))
		return
	}
	if author.FirstName == "" && author.LastName == "" {
		errorResponse(c, http.StatusNotFound, errors.New("author not found"))
		return
	}

	authorResponse := createResponseAuthors([]dataModels.Author{author})

	c.JSON(http.StatusOK, authorResponse)
}

func addAuthor(c *gin.Context) {
	var newAuthor dataModels.Author

	if err := c.ShouldBindJSON(&newAuthor); err != nil {
		errorResponse(c, http.StatusBadRequest, errors.New("there is something wrong with the data you sent"))
		logger.Error(err)
		return
	}
	if err := newAuthor.Validate(); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := db.AddAuthor(newAuthor); err != nil {
		if err.Error() == "author already exists" {
			errorResponse(c, http.StatusConflict, errors.New("author already exist"))
			return
		} else {
			errorResponse(c, http.StatusInternalServerError, errors.New("internal server error, please try making a request later"))
			logger.Error(err)
			return
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{"message": "author added successfully!"})
	}
}

func getByFullName(c *gin.Context, fullName string) {
	authors, err := db.SearchByFLName(fullName)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, errors.New("internal server error, please try making a request later"))
		logger.Error(err)
		return
	}
	if len(authors) == 0 {
		errorResponse(c, http.StatusNotFound, errors.New("record not found"))
		return
	}

	authorsResponse := createResponseAuthors(authors)

	c.JSON(http.StatusOK, authorsResponse)
}

func createResponseAuthors(authors []dataModels.Author) []dataModels.AuthorResponse {
	authorsResponse := make([]dataModels.AuthorResponse, len(authors))

	for i, author := range authors {
		authorsResponse[i] = dataModels.AuthorResponse{
			ID:        author.ID,
			FirstName: author.FirstName,
			LastName:  author.LastName,
		}
	}
	return authorsResponse
}
