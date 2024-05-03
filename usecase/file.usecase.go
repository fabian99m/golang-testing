package usecase

import (
	"bytes"
	"context"
	"dbtest/domain/dto"
	"dbtest/model"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type FileUseCase struct {
	bucketName string
	s3Client   S3API
}

type S3API interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func NewFileUseCase(bucketName string, s3API S3API) model.FileUseCase {
	if bucketName == "" {
		panic("bucketName requerido...")
	}
	return &FileUseCase{bucketName: bucketName, s3Client: s3API}
}

func (usecase *FileUseCase) GetFiles(token string) (*dto.FileResponseDto, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:  &usecase.bucketName,
		MaxKeys: aws.Int32(10),
	}

	if token == "" {
		input.ContinuationToken = nil
	} else {
		input.ContinuationToken = &token
	}

	output, listError := usecase.s3Client.ListObjectsV2(context.Background(), input)

	if listError != nil {
		log.Println("error obteniendo objects..", listError)
		return nil, listError
	}

	var keys []string
	for _, object := range output.Contents {
		keys = append(keys, *object.Key)
	}

	response := &dto.FileResponseDto{Keys: keys}
	if output.NextContinuationToken != nil {
		response.Next = *output.NextContinuationToken
	}

	return response, nil
}

func (usecase *FileUseCase) GetFile(key string) (*dto.FileResponseDto, error) {
	log.Println("Inicia descarga de archivo...")

	objectOutput, getError := usecase.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &usecase.bucketName,
		Key:    &key,
	})

	if getError != nil {
		log.Println("Error al descargar archivo...", getError)
		var NoSuchKey *types.NoSuchKey
		if errors.As(getError, &NoSuchKey) {
			return &dto.FileResponseDto{}, nil
		}

		return nil, getError
	}

	log.Println("Fin de descarga de archivo...")

	return dto.NewFileResponseDto(objectOutput.Body, *objectOutput.ContentType, *objectOutput.ContentLength), nil
}

func (usecase *FileUseCase) SaveFile(fileContent *[]byte, name *string) (string, error) {
	contentType := http.DetectContentType(*fileContent)
	log.Println("File Conten-Type:", contentType)

	fileName := generateFileName(name)
	log.Println("File name:", fileName)

	_, putError := usecase.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Body:        bytes.NewReader(*fileContent),
		Bucket:      aws.String(usecase.bucketName),
		Key:         &fileName,
		ContentType: &contentType,
	})

	if putError != nil {
		log.Println("Error guardando archivo en S3...", putError)
		return "", putError
	}

	return fileName, nil
}

func generateFileName(name *string) string {
	return *name + time.Now().Format("2006-01-02 15:04:05")
}
