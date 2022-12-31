package handler

import (
	"net/http"
	"strconv"

	"dbtest/domain/dto"
	"dbtest/model"

	"github.com/gin-gonic/gin"
)

type HeroHandler struct {
	usecase model.HeroUseCase
}

func NewHeroHandler(router *gin.RouterGroup, usecase model.HeroUseCase) {
	heroHandler := &HeroHandler{
		usecase: usecase,
	}

	router.GET("/heros", heroHandler.GetAll)
	router.GET("/heros/:id", heroHandler.GetHeroById)
	router.POST("/heros", heroHandler.SaveHero)
}

func (h *HeroHandler) GetAll(c *gin.Context) {
	response := h.usecase.GetAllHeros()
	c.JSON(response.Status, response)
}

func (h *HeroHandler) GetHeroById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseDto{Codigo: "1001", Mensaje: "Revisar formato de Id de entrada"})
		return
	}

	response := h.usecase.GetHeroById(id)

	c.JSON(response.Status, response)
}

func (h *HeroHandler) SaveHero(c *gin.Context) {
	var newHero dto.HeroDto

	if err := c.BindJSON(&newHero); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	response := h.usecase.SaveHero(newHero)
	c.JSON(response.Status, response)
}
