package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	HTTPCode int    `json:"code"`
	Error    string `json:"error,omitempty"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, ErrorResponse{
						HTTPCode: http.StatusInternalServerError,
						Error:    "panic recovered",
					})
				}
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()
			statusCode := c.Writer.Status()
			if statusCode == http.StatusOK {
				statusCode = http.StatusInternalServerError
			}

			

			log.Printf("DEBUG middleware: statusCode=%d, error=%q, written=%v",
				statusCode, lastErr.Error(), c.Writer.Written())

			if !c.Writer.Written() {
				c.JSON(statusCode, ErrorResponse{
					HTTPCode: statusCode,
					Error:    lastErr.Error(),
				})
			}
			c.Abort()
		}
	}
}
