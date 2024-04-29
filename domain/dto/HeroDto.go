package dto

type HeroDto struct {
	Id         int    `json:"id"`
	Name       string `json:"name" validate:"notblank"`
	CreateDate string `json:"create_date"`
}