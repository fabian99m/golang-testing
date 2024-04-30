package usecase

import (
	"bytes"
	"context"
	"dbtest/domain/dto"
	"dbtest/model"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileUseCase struct{
	bucketName string
	s3Client *s3.Client
}

func NewFileUseCase(bucketName string) model.FileUseCase { 
	if (bucketName == "") {
		panic("bucketName requerido...")
	}
	return &FileUseCase{bucketName: bucketName, s3Client: buildS3Client()}
}

func (usecase *FileUseCase) GetFile(key string) dto.FileResponseDto {
	log.Println("Inicia descarga de archivo...");

    objectOutput, objectOutputError := usecase.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(usecase.bucketName),
		Key: aws.String(key),
	})

	if objectOutputError != nil {
		log.Println("Error obteniendo objecto desde S3...", objectOutputError)
		return dto.FileResponseDto{Data: nil}
	}

	log.Println("Fin de descarga de archivo...")

	return dto.NewFileResponseDto(objectOutput.Body, *objectOutput.ContentType, *objectOutput.ContentLength);
}

func (usecase *FileUseCase) SaveFile(file *multipart.FileHeader) string {
	multipartFile, multipartFileError := file.Open()
	if multipartFileError != nil {
		log.Println("Error abriendo multipartFile...", multipartFileError)
		return "";
	}

	defer multipartFile.Close()

	fileContent, fileContentError := io.ReadAll(multipartFile)
	if fileContentError != nil {
		log.Println("Error leyendo todos los bytes de multipartFile", fileContentError)
		return "";
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

	return fileName;
}

func generateFileName(name *string) string {
	return *name + time.Now().Format("2006-01-02 15:04:05")
}

func buildS3Client() *s3.Client {
	cfg, cfgError := config.LoadDefaultConfig(context.Background())
	if cfgError != nil {
		log.Panic("unable to load SDK config,", cfgError)
	}

	return s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String("http://localhost:4566")
		options.UsePathStyle = true
	})
}