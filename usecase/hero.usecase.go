package usecase

import (
	"dbtest/model"
	"log"
	"net/http"
	"time"

	"dbtest/domain/dto"
	"dbtest/domain/mapper"
)

type HeroUseCase struct {
	repository model.HeroDbInteractor
}

func NewHeroUseCase(repository model.HeroDbInteractor) *HeroUseCase {
	return &HeroUseCase{repository: repository}
}

func (h *HeroUseCase) GetAllHeros() dto.ResponseDto {
	log.Println("GetAllHeros handler")

	heros, _ := h.repository.GetAll()

	return dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa", Data: mapper.ToHerosDto(heros)}
}

func (h *HeroUseCase) GetHeroById(id int) dto.ResponseDto {
	log.Println("GetHeroById handler")

	hero, RowsAffected := h.repository.GetById(id)
	if RowsAffected == 0 {
		return dto.ResponseDto{Status: http.StatusNotFound, Codigo: "1001", Mensaje: "Hero no encontrado"}
	}

	return dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa", Data: mapper.ToHeroDto(hero)}
}

func (h *HeroUseCase) SaveHero(newHero dto.HeroDto) dto.ResponseDto {
	log.Println("SaveHero handler")

	newHero.CreateDate = time.Now().Format("2006-01-02")

	hero := mapper.ToHero(newHero)

	if RowsAffected := h.repository.Save(&hero); RowsAffected == 0 {
		return dto.ResponseDto{Status: http.StatusInternalServerError, Codigo: "1002", Mensaje: "Error guardando hero"}
	}

	return dto.ResponseDto{
		Status:  http.StatusCreated,
		Codigo:  "1000",
		Mensaje: "Hero creado",
	}
}
