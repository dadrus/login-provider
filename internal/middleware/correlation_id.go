package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const headerName = "Correlation-Id"

func CorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.Request.Header.Get(headerName)
		if correlationId == "" {
			correlationId = uuid.New().String()
			c.Request.Header.Set(headerName, correlationId)
		}

		c.Next()

		c.Writer.Header().Set(headerName, correlationId)
	}
}
