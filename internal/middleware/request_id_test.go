package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIfRequestIdIsPresentItIsExtended(t *testing.T) {
	// GIVEN
	requestId := "foo:bar"
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{ Header: make(http.Header)}
	ctx.Request.Header.Set(requestIdHeaderName, requestId)
	middleware := RequestId()

	// WHEN
	middleware(ctx)

	// THEN
	val := ctx.Request.Header.Get(requestIdHeaderName)
	require.NotEmpty(t, val, "X-Request-Id header must not be empty")

	reqIdComponents := strings.Split(val, ";")
	require.Equal(t, len(reqIdComponents), 2,"X-Request-Id must consist of two components")
	require.Equal(t, reqIdComponents[0], requestId, "First X-Request-Id component must be equal the value of the sent X-Request-Id")

	setReqId := strings.Split(reqIdComponents[1], ":")
	require.Equal(t, len(setReqId), 2,"Set X-Request-Id component must consist of two parts")
	require.Equal(t, setReqId[0], "login-provider", "First part of the set X-Request-Id component must be 'login-provider'")
	require.NotEmpty(t, setReqId[1], "Second part of the set X-Request-Id component must not be empty")
}

func TestIfRequestIdIsNotPresentItIsAdded(t *testing.T) {
	// GIVEN
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{ Header: make(http.Header)}
	middleware := RequestId()

	// WHEN
	middleware(ctx)

	// THEN
	val := ctx.Request.Header.Get(requestIdHeaderName)
	require.NotEmpty(t, val, "X-Request-Id header must not be empty")

	reqIdComponents := strings.Split(val, ";")
	require.Equal(t, len(reqIdComponents), 1,"X-Request-Id must consist of one component")

	setReqId := strings.Split(reqIdComponents[0], ":")
	require.Equal(t, len(setReqId), 2,"Set X-Request-Id component must consist of two parts")
	require.Equal(t, setReqId[0], "login-provider", "First part of the set X-Request-Id component must be 'login-provider'")
	require.NotEmpty(t, setReqId[1], "Second part of the set X-Request-Id component must not be empty")
}

func TestRequestIdMiddlewareCallsContextNext(t *testing.T) {
	// TODO: How to achieve that. gin.Context is a struct and as such not mockable
}
