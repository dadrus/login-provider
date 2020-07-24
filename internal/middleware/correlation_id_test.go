package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestIfCorrelationIdIsPresentItIsReused(t *testing.T) {
	// GIVEN
	correlationId := "foo"
	ctx := gin.Context{
		Request: &http.Request{ Header: make(http.Header)},
		Writer:  &MockedResponseWriter{},
	}
	ctx.Request.Header.Set(correlationIdHeaderName, correlationId)
	middleware := CorrelationId()

	// WHEN
	middleware(&ctx)

	// THEN
	val := ctx.Writer.Header().Get(correlationIdHeaderName)
	require.NotEmpty(t, val, "Correlation-Id header must not be empty")
	require.Equal(t, correlationId, val, "Set and received correlation ids must be equal")
}

func TestIfCorrelationIdIsNotPresentItIsCreated(t *testing.T) {
	// GIVEN
	ctx := gin.Context{
		Request: &http.Request{ Header: make(http.Header)},
		Writer:  &MockedResponseWriter{},
	}
	middleware := CorrelationId()

	// WHEN
	middleware(&ctx)

	// THEN
	reqVal := ctx.Request.Header.Get(correlationIdHeaderName)
	respVal := ctx.Writer.Header().Get(correlationIdHeaderName)
	require.NotEmpty(t, reqVal, "Correlation-Id header must not be empty")
	require.Equal(t, reqVal, respVal, "Request and response correlation ids must be equal")
}

func TestCorrelationIdMiddlewareCallsContextNext(t *testing.T) {
	// TODO: How to achieve that. gin.Context is a struct and as such not mockable
}
