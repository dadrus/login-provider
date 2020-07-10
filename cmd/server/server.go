package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"login-provider/internal/config"
	"login-provider/internal/handler"
	"login-provider/internal/hydra"
	"os"
	"time"
)

func Serve(cmd *cobra.Command, args []string) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(correlationIdFilter())
	router.Use(logger())
	router.LoadHTMLGlob("web/templates/*")

	// TODO: refactor this part and move it closer to the handler functions
	hf, err := hydra.NewClientFactory(viper.GetString(config.HydraAdminUrl))
	if err != nil {
		l := log.With().Err(err).Logger()
		l.Fatal().Msg("Failed to create hydra client factory")
	}

	router.GET("/login", handler.ShowLoginPage(hf))
	router.POST("/login", handler.Login(hf))
	router.GET("/consent", handler.ShowConsentPage(hf))
	router.POST("/consent", handler.Consent(hf))
	router.GET("/logout", handler.ShowLogoutPage(hf))
	router.POST("/logout", handler.Logout(hf))

	addr := configuredAddress()

	if tlsConfig, err := getTlsConfig(); err == nil {
		log.Info().
			Msg("Listening and serving HTTPS on " + addr)
		router.RunTLS(addr, tlsConfig.TlsCertFile, tlsConfig.TlsKeyFile)
	} else {
		log.Info().
			Msg("Listening and serving HTTP on " + addr)
		router.Run(addr)
	}
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		l := log.Logger.With().
			Str("_ops_correlation_id", c.Request.Header.Get("Correlation-Id")).
			Logger()

		c.Set("logger", l)

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
			Str("_http_x_request_id", c.Request.Header.Get("X-Request-Id")).
			Str("_http_x_amz_cf_id", c.Request.Header.Get("X-Amz-Cf-Id")).
			Int64("_ops_tx_start", start.Unix()).
			Dur("_opx_tx_duration", latency).
			Msg("tx")
	}
}

func correlationIdFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.Request.Header.Get("Correlation-Id")
		if correlationId == "" {
			correlationId = uuid.New().String()
			c.Request.Header.Set("Correlation-Id", correlationId)
		}

		c.Next()

		c.Writer.Header().Set("Correlation-Id", correlationId)
	}
}

type tlsConfig struct {
	TlsKeyFile  string
	TlsCertFile string
}

func getTlsConfig() (*tlsConfig, error) {
	tlsKeyFile := viper.GetString(config.TlsKeyFile)
	if len(tlsKeyFile) == 0 {
		return nil, errors.New("no TLS key configured")
	}
	if _, err := os.Stat(tlsKeyFile); err != nil {
		return nil, errors.New("configured TLS key not available")
	}

	tlsCertFile := viper.GetString(config.TlsCertFile)
	if len(tlsCertFile) == 0 {
		return nil, errors.New("no TLS cert configured")
	}
	if _, err := os.Stat(tlsCertFile); err != nil {
		return nil, errors.New("configured TLS cert not available")
	}

	return &tlsConfig{
		TlsKeyFile:  tlsKeyFile,
		TlsCertFile: tlsCertFile,
	}, nil
}

func configuredAddress() string {
	return viper.GetString(config.Host) + ":" + viper.GetString(config.Port)
}
