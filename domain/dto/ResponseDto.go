package dto

type ResponseDto struct {
	Status  int    `json:"-"`
	Codigo  string `json:"codigo"`
	Mensaje string `json:"mensaje"`
	Data    any    `json:"data"`
}
