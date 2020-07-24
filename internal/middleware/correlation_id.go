package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const correlationIdHeaderName = "Correlation-Id"

func CorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.Request.Header.Get(correlationIdHeaderName)
		if correlationId == "" {
			correlationId = uuid.New().String()
			c.Request.Header.Set(correlationIdHeaderName, correlationId)
		}

		c.Next()

		c.Writer.Header().Set(correlationIdHeaderName, correlationId)
	}
}
