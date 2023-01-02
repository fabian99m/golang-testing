package main

import (
	"dbtest/db"
	"dbtest/handler"
	"dbtest/model"
	"dbtest/repository"
	"dbtest/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

func initApp() model.HeroUseCase {
	dbConnection := db.NewDbConnection()
	heroRepository := repository.NewHeroRespository(dbConnection)
	return usecase.NewHeroUseCase(heroRepository)
}

func main() {
	useCase := initApp()
	router := gin.Default()
	api := router.Group("/v1")

	handler.NewHeroHandler(api, useCase)

	log.Panic(router.Run("localhost:3000"))
}
