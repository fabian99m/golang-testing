package dto

import "fmt"

type ResponseDto struct {
	Status  int    `json:"-"`
	Codigo  string `json:"codigo"`
	Mensaje string `json:"mensaje"`
	Data    any    `json:"data"`
}



func (response ResponseDto) Error() string {
    return fmt.Sprintf("codigo: %s,  metadata: %s", response.Codigo, response.Mensaje)
}

var InvalidData = ResponseDto{Codigo: "5501", Status: 400, Mensaje: "Data invalida"};