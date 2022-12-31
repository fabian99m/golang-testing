package model

import (
	"dbtest/domain/dto"
	"time"
)

type Hero struct {
	Id         int       `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(255)"`
	CreateDate time.Time `gorm:"type:date"`
}

type HeroDbInteractor interface {
	GetAll(dest *[]Hero)
	GetById(id int, dest *Hero) int64
	Save(dest *Hero) int64
}

type HeroUseCase interface {
	GetAllHeros() dto.ResponseDto
	GetHeroById(id int) dto.ResponseDto
	SaveHero(hero dto.HeroDto) dto.ResponseDto
}
