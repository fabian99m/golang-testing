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

func (hr *heroRespository) GetAll(dest *[]model.Hero) {
	hr.bd.Find(dest)
}

func (hr *heroRespository) GetById(id int, dest *model.Hero) int64 {
	r := hr.bd.First(dest, id)
	return r.RowsAffected
}

func (hr *heroRespository) Save(dest *model.Hero) int64 {
	r := hr.bd.Create(dest)
	return r.RowsAffected
}
