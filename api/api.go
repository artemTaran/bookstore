package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"shop/db"
	"shop/logger"
)

func InitRouter() {
	router := gin.Default()

	router.GET("/api/books", func(c *gin.Context) {
		c.Request.URL.Path = "/api/books/10"
		router.HandleContext(c)
	})
	router.GET("/api/books/:quantity", getBooks)
	router.GET("/api/books/byid/:id", getBookById)
	router.POST("/api/books", addBook)
	router.GET("/api/books/search/:title", searchByTitle)

	router.GET("/api/authors", func(c *gin.Context) {
		c.Request.URL.Path = "/api/authors/10"
		router.HandleContext(c)
	})
	router.GET("/api/authors/:quantity", getAuthors)
	router.GET("/api/authors/byid/:id", getAuthorById)
	router.POST("/api/authors", addAuthor)
	router.GET("/api/authors/search/:flName", searchByFLName)

	host, port := getEnvV()
	address := fmt.Sprintf("%s:%s", host, port)

	db.InitDb()

	err := router.Run(address)
	if err != nil {
		logger.Fatal(err)
		return
	}
}

func getEnvV() (host, port string) {
	host = os.Getenv("HTTP_HOST")
	port = os.Getenv("HTTP_PORT")
	if host == "" {
		logger.Fatal(fmt.Errorf("HTTP_HOST is not set"))
		return
	}
	if port == "" {
		logger.Fatal(fmt.Errorf("HTTP_PORT is not set"))
		return
	}
	return host, port
}
