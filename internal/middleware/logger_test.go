package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type stringWriter struct {
	value string
}

// The only function required by the io.Writer interface.  Will append
// written data to the String.value string.
func (s *stringWriter) Write(p []byte) (n int, err error) {
	s.value += string(p)
	return len(p), nil
}

func TestLoggerMiddlewareAddsRequiredMDC(t *testing.T) {
	// GIVEN
	requestId := "foo:bar"
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{ Header: make(http.Header)}
	ctx.Request.Header.Set(requestIdHeaderName, requestId)
	ctx.Request.URL = &url.URL{
		Path: "foo",
		RawQuery: "bar",
	} 
	middleware := Logger()
	w := &stringWriter{}
	log.Logger = log.Output(w)

	// WHEN
	middleware(ctx)

	// THEN
	assert.Contains(t, w.value, "level")
	assert.Contains(t, w.value, "_ops_correlation_id")
	assert.Contains(t, w.value, "_http_x_request_id")
	assert.Contains(t, w.value, "_ops_caller")
	assert.Contains(t, w.value, "_ops_tx_method")
	assert.Contains(t, w.value, "_ops_tx_object")
	assert.Contains(t, w.value, "_ops_tx_result_code")
	assert.Contains(t, w.value, "_ops_tx_body_bytes_sent")
	assert.Contains(t, w.value, "_ops_tx_scheme")
	assert.Contains(t, w.value, "_http_x_forwarded_host")
	assert.Contains(t, w.value, "_http_x_forwarded_for")
	assert.Contains(t, w.value, "_http_x_forwarded_port")
	assert.Contains(t, w.value, "_http_x_forwarded_proto")
	assert.Contains(t, w.value, "_http_user_agent")
	assert.Contains(t, w.value, "_http_x_request_id")
	assert.Contains(t, w.value, "_http_x_amz_cf_id")
	assert.Contains(t, w.value, "_ops_tx_start")
	assert.Contains(t, w.value, "_opx_tx_duration")
	assert.Contains(t, w.value, "time")
	assert.Contains(t, w.value, "\"message\":\"tx\"")
}
