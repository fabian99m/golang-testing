package model

import (
	"dbtest/domain/dto"
)

type FileUseCase interface {
	SaveFile(*[]byte, *string) (string, error)
	GetFile(string) (*dto.FileResponseDto, error)
	GetFiles(string) (*dto.FileResponseDto, error)
}
