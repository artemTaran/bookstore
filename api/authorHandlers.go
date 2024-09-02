package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop/dataModels"
	"shop/db"
	"strconv"
)

func getAuthors(c *gin.Context) {
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
	authors, err := db.GetAuthors(quantity)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
	}

	authorsResponse := createResponseAuthors(authors)

	c.JSON(http.StatusOK, authorsResponse)
}

func getAuthorById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid id"))
		return
	}
	author, err := db.GetAuthorById(id)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("author not found"))
		return
	}
	if author.FirstName == "" && author.LastName == "" {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("author not found"))
		return
	}

	authorResponse := createResponseAuthors([]dataModels.Author{author})

	c.JSON(http.StatusOK, authorResponse)
}

func addAuthor(c *gin.Context) {
	var newAuthor dataModels.Author

	if err := c.ShouldBindJSON(&newAuthor); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	if newAuthor.FirstName == "" || newAuthor.LastName == "" {
		errorResponse(c, http.StatusBadRequest, fmt.Errorf("missing required fields"))
		return
	}

	msg, err := db.AddAuthor(newAuthor)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": msg})
}

func searchByFLName(c *gin.Context) {
	flName := c.Param("flName")
	authors, err := db.SearchByFLName(flName)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return
	}
	if len(authors) == 0 {
		errorResponse(c, http.StatusNotFound, fmt.Errorf("record not found"))
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
