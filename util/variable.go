package util

import (
	"github.com/joho/godotenv"
	"os"
)

func GetVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	return os.Getenv(key)
}
