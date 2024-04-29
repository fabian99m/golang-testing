package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
) 

func S3connect2(c *gin.Context) {

	multipart, errorFile := c.FormFile("file")

	if errorFile != nil {
		c.Status(200)
		return;
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Println("unable to load SDK config,", err)
		return
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String("http://localhost:4566")
	})

	file, _ := multipart.Open();

	resp, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Body: file,
		Bucket: aws.String("testbucket"),
		Key:  aws.String(time.Now().Format("2006-01-02 15:04:05")),
		
	})
	c.JSON(200, resp.ResultMetadata)
}