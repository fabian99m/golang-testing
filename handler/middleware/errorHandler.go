package middleware

import (
	h "dbtest/handler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch err.Err {
		case h.ErrNotFound:
			log.Println("h.ErrNotFound")
			c.JSON(-1, gin.H{"error": h.ErrNotFound.Error()})
		case h.ErrBadRequest:
			log.Println("h.ErrBadRequest")
			c.JSON(-1, gin.H{"error en entrada": h.ErrNotFound.Error()})
		default:
			log.Println("default error")
			c.JSON(-1, gin.H{"error inesperado": "test"})
		}
	}

	c.JSON(http.StatusInternalServerError, "")
}
