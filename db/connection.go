package db

import (
	"dbtest/config"
	"dbtest/model"

	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewDbConnection(dataBaseConfig *config.DataBase) *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dataBaseConfig.Host, dataBaseConfig.Port, dataBaseConfig.UserName, dataBaseConfig.Password, dataBaseConfig.Name)

	connection, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: getLogger(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dataBaseConfig.Esquema + ".",
			SingularTable: true,
		},
	})

	if err != nil {
		panic("Failed to connect database")
	}

	if err = connection.AutoMigrate(&model.Hero{}); err != nil {
		panic("Failed database migration")
	}

	return connection
}

func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
