package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIfCorrelationIdIsPresentItIsReused(t *testing.T) {
	// GIVEN
	correlationId := "foo"
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{ Header: make(http.Header)}
	ctx.Request.Header.Set(correlationIdHeaderName, correlationId)
	middleware := CorrelationId()

	// WHEN
	middleware(ctx)

	// THEN
	val := ctx.Writer.Header().Get(correlationIdHeaderName)
	require.NotEmpty(t, val, "Correlation-Id header must not be empty")
	require.Equal(t, correlationId, val, "Set and received correlation ids must be equal")
}

func TestIfCorrelationIdIsNotPresentItIsCreated(t *testing.T) {
	// GIVEN
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{ Header: make(http.Header)}
	middleware := CorrelationId()

	// WHEN
	middleware(ctx)

	// THEN
	reqVal := ctx.Request.Header.Get(correlationIdHeaderName)
	respVal := ctx.Writer.Header().Get(correlationIdHeaderName)
	require.NotEmpty(t, reqVal, "Correlation-Id header must not be empty")
	require.Equal(t, reqVal, respVal, "Request and response correlation ids must be equal")
}
