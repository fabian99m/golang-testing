package main

import (
	"dbtest/config"
	"dbtest/handler"
	"dbtest/handler/middleware"

	"log"
	"fmt"
	

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler)
	api := router.Group("/v1")

	handler.NewHandler(api, readConfig())

	log.Panic(router.Run("localhost:3000"))
}

func readConfig() *config.Config {
	cfg, errCfg := config.ReadConf()

	if errCfg != nil {
		fmt.Println("Error en yml...", errCfg)
		panic("Failed to load config from YML...")
	}

	return cfg
}
