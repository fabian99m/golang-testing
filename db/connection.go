package db

import (
	"dbtest/model"
	"dbtest/util"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Connection *gorm.DB

var (
	host      = util.GetVariable("BD_HOST")
	port      = util.GetVariable("BD_PORT")
	user      = util.GetVariable("BD_USERNAME")
	password  = util.GetVariable("BD_PASSWORD")
	dbname    = util.GetVariable("BD_NAME")
	bd_schema = util.GetVariable("BD_ESQUEMA")+"."
)

func Connect() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	Connection, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: getLogger(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: bd_schema,
			SingularTable: true,
		},
	})

	Connection.AutoMigrate(&model.Hero{})

	if err != nil {
		panic("Failed to connect database")
	}
}

func getLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Info,
			IgnoreRecordNotFoundError: true, 
			Colorful: true,
		},
	)
}