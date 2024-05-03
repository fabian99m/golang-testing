package app

import (
	"dbtest/config"
	"dbtest/db"
	"dbtest/model"
	"dbtest/repository"
	"dbtest/usecase"

	"context"
	"log"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type AppConfig struct {
	HeroUseCase model.HeroUseCase
	FileUseCase model.FileUseCase
}

func NewAppConfig() *AppConfig {
	config := readConfig()
	dbConnection := db.NewDbConnection(&config.DataBase)
	heroRepository := repository.NewHeroRespository(dbConnection)

	return &AppConfig{
		HeroUseCase: usecase.NewHeroUseCase(heroRepository),
		FileUseCase: usecase.NewFileUseCase(config.S3.Bucket, buildS3Client()),
	}
}

func readConfig() *config.Config {
	cfg, errCfg := config.ReadConf()

	if errCfg != nil {
		log.Println("Error en yml...", errCfg)
		panic("Failed to load config from YML...")
	}

	return cfg
}

func buildS3Client() *s3.Client {
	cfg, cfgError := awsconfig.LoadDefaultConfig(context.Background())
	if cfgError != nil {
		log.Panic("unable to load SDK config,", cfgError)
	}

	return s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String("http://localhost:4566")
		options.UsePathStyle = true
	})
}
