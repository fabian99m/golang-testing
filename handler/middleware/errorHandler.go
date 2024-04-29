package middleware

import (
	"dbtest/domain/dto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch e := err.Err.(type) {
		case dto.ResponseDto:
			c.AbortWithStatusJSON(e.Status, e)
		default:
			log.Println("default error", e)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error inesperado": "test"})
		}
	}
}