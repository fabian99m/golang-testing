package usecase

import (
	"dbtest/model"
	"log"
	"net/http"

	"dbtest/domain/dto"
	"dbtest/domain/mapper"
)

type HeroUseCase struct {
	repository model.HeroDbInteractor
}

func NewHeroUseCase(repository model.HeroDbInteractor) model.HeroUseCase {
	return &HeroUseCase{repository: repository}
}

func (h *HeroUseCase) GetAllHeros() dto.ResponseDto {
	log.Println("GetAllHeros handler")

	var heros []model.Hero
	h.repository.GetAll(&heros)

	return dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa", Data: mapper.ToHerosDto(heros)}
}

func (h *HeroUseCase) GetHeroById(id int) dto.ResponseDto {
	log.Println("GetHeroById handler")

	var hero model.Hero
	if RowsAffected := h.repository.GetById(id, &hero); RowsAffected == 0 {
		return dto.ResponseDto{Status: http.StatusNotFound, Codigo: "1001", Mensaje: "Hero no encontrado"}
	}

	return dto.ResponseDto{Status: http.StatusOK, Codigo: "1000", Mensaje: "Consulta exitosa", Data: mapper.ToHeroDto(hero)}
}

func (h *HeroUseCase) SaveHero(newHero dto.HeroDto) dto.ResponseDto {
	log.Println("SaveHero handler")

	hero, err := mapper.ToHero(newHero)

	if err != nil {
		return dto.ResponseDto{Status: http.StatusInternalServerError, Codigo: "1003", Mensaje: "Error convirtiendo fecha"}
	}

	if RowsAffected := h.repository.Save(&hero); RowsAffected == 0 {
		return dto.ResponseDto{Status: http.StatusInternalServerError, Codigo: "1002", Mensaje: "Error guardando hero"}
	}

	return dto.ResponseDto{
		Status:  http.StatusCreated,
		Codigo:  "1000",
		Mensaje: "Hero creado",
	}
}
