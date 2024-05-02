package model

import (
	"mime/multipart"
	"dbtest/domain/dto"
)

type FileUseCase interface {
	SaveFile(*multipart.FileHeader) string;
	GetFile(string) (*dto.FileResponseDto, error);
	GetFiles(string) (*dto.FileResponseDto, error);
}