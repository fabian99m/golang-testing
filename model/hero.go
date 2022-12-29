
package model

import "time"

type Hero struct {
	Id int `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255)"`
	CreateDate time.Time `gorm:"type:date"`
}