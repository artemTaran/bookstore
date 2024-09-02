package main

import (
	"shop/api"
	"shop/logger"
)

func main() {
	defer logger.Sync()
	api.InitRouter()
}
