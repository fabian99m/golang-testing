package mapper

import (
	"dbtest/domain/dto"
	"dbtest/model"
	"time"
)

func ToHeroDto(model model.Hero) dto.HeroDto {
	return dto.HeroDto{
		Id:         model.Id,
		Name:       model.Name,
		CreateDate: model.CreateDate.Format("2006-01-02"),
	}
}

func ToHerosDto(models []model.Hero) []dto.HeroDto {
	dtos := make([]dto.HeroDto, len(models))
	for i, m := range models {
		dtos[i] = ToHeroDto(m)
	}
	return dtos
}

func ToHero(dto dto.HeroDto) (model.Hero, error) {
	date := time.Now()

	return model.Hero{Name: dto.Name, CreateDate: date}, nil
}
