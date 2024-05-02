package usecase

import (
	"bytes"
	"context"
	"dbtest/domain/dto"
	"dbtest/model"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type FileUseCase struct {
	bucketName string
	s3Client   *s3.Client
}

func NewFileUseCase(bucketName string, s3Client  *s3.Client) model.FileUseCase {
	if bucketName == "" {
		panic("bucketName requerido...")
	}
	return &FileUseCase{bucketName: bucketName, s3Client: s3Client}
}

func (usecase *FileUseCase) GetFiles(token string) (*dto.FileResponseDto, error) {
	params := &s3.ListObjectsV2Input{
		Bucket: &usecase.bucketName,
		MaxKeys: aws.Int32(10),
	};

	if token == "" { params.ContinuationToken = nil} else { params.ContinuationToken = &token }

	objects, listError := usecase.s3Client.ListObjectsV2(context.Background(), params)

	if listError != nil {
		log.Println("error obteniendo objects..", listError)
		return nil, listError;
	}

	var keys []string
	for _, object := range objects.Contents {
		keys = append(keys, *object.Key)
	}

	response := &dto.FileResponseDto{Keys: keys}
	if (objects.NextContinuationToken != nil) {
		response.Next = *objects.NextContinuationToken;
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
		var apiErr smithy.APIError
		if errors.As(getError, &apiErr) {
			switch apiErr.(type) {
			case *types.NoSuchKey:
				log.Println("Error NoSuchKey...")
				return &dto.FileResponseDto{}, nil
			case *types.NoSuchBucket:
				log.Println("Error NoSuchBucket...")
			default:
				log.Println("Error default...", getError)
				return &dto.FileResponseDto{}, getError
			}
		}
	}

	log.Println("Fin de descarga de archivo...")

	return dto.NewFileResponseDto(objectOutput.Body, *objectOutput.ContentType, *objectOutput.ContentLength), nil
}

func (usecase *FileUseCase) SaveFile(file *multipart.FileHeader) string {
	multipartFile, multipartFileError := file.Open()
	if multipartFileError != nil {
		log.Println("Error abriendo multipartFile...", multipartFileError)
		return ""
	}

	defer multipartFile.Close()

	fileContent, fileContentError := io.ReadAll(multipartFile)
	if fileContentError != nil {
		log.Println("Error leyendo todos los bytes de multipartFile", fileContentError)
		return ""
	}

	contentType := http.DetectContentType(fileContent)
	log.Println("File Conten-Type:", contentType)

	fileName := generateFileName(&file.Filename)
	log.Println("File name:", fileName)

	_, putError := usecase.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Body:        bytes.NewReader(fileContent),
		Bucket:      aws.String(usecase.bucketName),
		Key:         &fileName,
		ContentType: &contentType,
	})
	if putError != nil {
		log.Println("Error guardando archivo en S3...", putError)
		return ""
	}

	return fileName
}

func generateFileName(name *string) string {
	return *name + time.Now().Format("2006-01-02 15:04:05")
}