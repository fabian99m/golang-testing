package dto

type ResponseDto struct {
	Codigo string `json:"codigo"`
	Mensaje string `json:"mensaje"`
	Data any `json:"data"`
}