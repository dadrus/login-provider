package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type MockedResponseWriter struct {
	gin.ResponseWriter
	header http.Header
}

func (rw *MockedResponseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = make(http.Header)
	}

	return rw.header
}
