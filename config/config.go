package config

import (
	util "dbtest/util"
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DataBase DataBase `yaml:"database"`
	S3 S3 `yaml:"s3"`
}

type S3 struct {
	Bucket string `yaml:"bucket" validate:"notblank"`
}

type DataBase struct {
	Host     string `yaml:"host" validate:"notblank"`
	Port     int32  `yaml:"port" validate:"min=1"`
	UserName string `yaml:"username" validate:"notblank"`
	Password string `yaml:"password" validate:"notblank"`
	Name     string `yaml:"name" validate:"notblank"`
	Esquema  string `yaml:"esquema" validate:"notblank"`
}

func ReadConf() (*Config, error) {
	buffer, errRead := ioutil.ReadFile("config.yml")
	if errRead != nil {
		return nil, errRead
	}

	config := &Config{}
	errRead = yaml.Unmarshal(buffer, config)
	if errRead != nil {
		return nil, errRead
	}

	valid, fields := util.Validate(config); if !valid {
		log.Println("los siguiente campos son invalidos: ", fields)
		return nil, errors.New("broken")
	}

	return config, nil
}