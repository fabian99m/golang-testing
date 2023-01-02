package repository

import (
	"dbtest/model"
	"gorm.io/gorm"
)

type heroRespository struct {
	bd *gorm.DB
}

func NewHeroRespository(bd *gorm.DB) model.HeroDbInteractor {
	return &heroRespository{bd: bd}
}

func (hr *heroRespository) GetAll() ([]model.Hero, int64) {
	var dest []model.Hero
	r := hr.bd.Find(&dest)
	return dest, r.RowsAffected
}

func (hr *heroRespository) GetById(id int) (model.Hero, int64) {
	var dest model.Hero
	r := hr.bd.First(&dest, id)
	return dest, r.RowsAffected
}

func (hr *heroRespository) Save(hero *model.Hero) int64 {
	r := hr.bd.Create(hero)
	return r.RowsAffected
}
