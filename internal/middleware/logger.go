package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		l := log.With().
			Str("_ops_correlation_id", c.Request.Header.Get("Correlation-Id")).
			Str("_http_x_request_id", c.Request.Header.Get("X-Request-Id")).
			Logger()

		newCtx := l.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(newCtx)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		end = end.UTC()

		l.Info().
			Str("_ops_caller", c.Request.RemoteAddr).
			Str("_ops_tx_method", c.Request.Method).
			Str("_ops_tx_object", path).
			Int("_ops_tx_result_code", c.Writer.Status()).
			Int("_ops_tx_body_bytes_sent", c.Writer.Size()).
			Str("_ops_tx_scheme", c.Request.URL.Scheme).
			Str("_http_x_forwarded_host", c.Request.Header.Get("X-Forwarded-Host")).
			Str("_http_x_forwarded_for", c.Request.Header.Get("X-Forwarded-For")).
			Str("_http_x_forwarded_port", c.Request.Header.Get("X-Forwarded-Port")).
			Str("_http_x_forwarded_proto", c.Request.Header.Get("X-Forwarded-Proto")).
			Str("_http_user_agent", c.Request.UserAgent()).
			Str("_http_x_amz_cf_id", c.Request.Header.Get("X-Amz-Cf-Id")).
			Int64("_ops_tx_start", start.Unix()).
			Dur("_opx_tx_duration", latency).
			Msg("tx")
	}
}
