package main

import (
	"dbtest/db"
	"dbtest/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db.Connect()

	r.GET("/heros", handler.GetAllHeros)
	r.GET("/heros/:id", handler.GetHeroById)
	r.POST("/heros", handler.SaveHero)

	r.Run("localhost:3000")
}