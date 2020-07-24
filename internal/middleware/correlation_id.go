package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.Request.Header.Get("Correlation-Id")
		if correlationId == "" {
			correlationId = uuid.New().String()
			c.Request.Header.Set("Correlation-Id", correlationId)
		}

		c.Next()

		c.Writer.Header().Set("Correlation-Id", correlationId)
	}
}
