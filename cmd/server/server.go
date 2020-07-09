package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"login-provider/internal/config"
	"login-provider/internal/handler"
	"login-provider/internal/hydra"
	"net/http"
	"os"
	"time"
)

func Serve(cmd *cobra.Command, args []string) {
	router := gin.New()
	router.Use(logger())
	router.Use(gin.Recovery())
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

	if tlsConfig, err := getTlsConfig(); err == nil {
		log.Info().
			Msg("Listening and serving HTTPS")
		router.RunTLS(configuredAddress(), tlsConfig.TlsCertFile, tlsConfig.TlsKeyFile)
	} else {
		log.Info().
			Msg("Listening and serving HTTP")
		router.Run(configuredAddress())
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
			Str("host", c.Request.Host).
			Str("_method", c.Request.Method).
			Str("_path", path).
			Str("_user_agent", c.Request.UserAgent()).
			Logger()

		c.Set("logger", l)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		end = end.UTC()

		msg := "Request"
		if len(c.Errors) > 0 {
			msg = c.Errors.String()
		}

		l = l.With().
			Int("_ops_tx_result_code", c.Writer.Status()).
			Dur("_opx_tx_duration", latency).
			Logger()

		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			l.Warn().Msg(msg)
		case c.Writer.Status() >= http.StatusInternalServerError:
			l.Error().Msg(msg)
		default:
			l.Info().Msg(msg)
		}

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
