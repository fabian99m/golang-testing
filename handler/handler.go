package handler

import (
	"dbtest/model"
	"dbtest/repository"
	"dbtest/domain/mapper"
	"dbtest/domain/dto"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllHeros(c *gin.Context) {
	log.Println("GetAllHeros handler")
	var heros []model.Hero

	repository.GetAll(&heros)

	c.JSON(http.StatusOK, dto.ResponseDto{Codigo: "1000", Mensaje: "Consulta exitosa", Data: mapper.ToHerosDto(heros)})
}

func GetHeroById(c *gin.Context) {
	log.Println("GetHeroById handler")

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseDto{Codigo: "1001", Mensaje: "Revisar formato de Id de entrada",})
		return
	}

	var hero model.Hero
	if result := repository.GetById(id, &hero); result.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, dto.ResponseDto{Codigo: "1000", Mensaje: "Consulta exitosa", Data: mapper.ToHeroDto(hero)})
}

func SaveHero(c *gin.Context) {
	log.Println("SaveHero handler")

	var newHero dto.HeroDto
	if err := c.BindJSON(&newHero); err != nil {
		log.Panic(err)
		return
	}

	hero, err := mapper.ToHero(newHero)

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseDto{ Codigo: "1001", Mensaje: "Error convirtiendo fecha", })
		return
	}

	repository.Save(&hero)

	c.JSON(http.StatusCreated, dto.ResponseDto{Codigo: "1000", Mensaje: "Hero creado",})
}