package main

import (
	"dbtest/handler"
	"dbtest/app"
	"dbtest/handler/middleware"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler)
	api := router.Group("/v1")

	handler.NewHandler(api, app.NewAppConfig())

	log.Panic(router.Run("localhost:3000"))
}