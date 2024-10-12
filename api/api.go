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

	router.POST("/api/books", addBook)
	router.GET("/api/books", getBooks)

	router.GET("/api/authors", getAuthors)
	router.POST("/api/authors", addAuthor)

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
