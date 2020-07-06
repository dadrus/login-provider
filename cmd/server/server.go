package server

import (
	"errors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"login-provider/internal/config"
	"login-provider/internal/view"
	"os"
)

func Serve(cmd *cobra.Command, args []string) {
	router := gin.New()
	router.Use(logger.SetLogger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("web/templates/*")

	// TODO: refactor this part and move it closer to the handler functions
	hf := view.NewHydraClientFactory(viper.GetString(config.HydraAdminUrl))

	router.GET("/login", view.ShowLoginPage(hf))
	router.POST("/login", view.Login(hf))
	router.GET("/consent", view.ShowConsentPage(hf))
	router.POST("/consent", view.Consent(hf))
	router.GET("/logout", view.ShowLogoutPage(hf))
	router.POST("/logout", view.Logout(hf))

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
