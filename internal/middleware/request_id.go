package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		gotRequestId := c.Request.Header.Get("X-Request-Id")
		requestId := "login-provider:" + uuid.New().String()
		if len(gotRequestId) != 0 {
			requestId = gotRequestId + ";" + requestId
		}
		c.Request.Header.Set("X-Request-Id", requestId)

		c.Next()
	}
}
