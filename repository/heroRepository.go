package repository

import (
	"dbtest/db"
	"dbtest/model"

	"gorm.io/gorm"
)

func GetAll(dest *[]model.Hero) *gorm.DB {
	return db.Connection.Find(dest)
}

func GetById(id int, dest *model.Hero) *gorm.DB {
	return db.Connection.First(dest, id)
}

func Save(dest *model.Hero) *gorm.DB {
	return db.Connection.Create(dest)
}