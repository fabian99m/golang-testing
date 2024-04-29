package dto

import "io"

type FileResponseDto struct {
	Data          io.ReadCloser
	ContentType   string
	ContentLength int64
}

func NewFileResponseDto(data io.ReadCloser, contenType string, contentLength int64) FileResponseDto {
	return FileResponseDto{
		Data: data,
		ContentType: contenType,
		ContentLength: contentLength,
	};
}