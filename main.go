package main

import (
	"fmt"
	"shop/api"
	"shop/logger"
)

func main() {

	defer logger.Sync()
	logger.Error(fmt.Errorf("ErrorTest"))
	logger.Info("InfoTest")
	api.InitRouter()
	fmt.Println("d")

}
