package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIdHeaderName = "X-Request-Id"

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		gotRequestId := c.Request.Header.Get(requestIdHeaderName)
		requestId := "login-provider:" + uuid.New().String()
		if len(gotRequestId) != 0 {
			requestId = gotRequestId + ";" + requestId
		}
		c.Request.Header.Set(requestIdHeaderName, requestId)

		c.Next()
	}
}
