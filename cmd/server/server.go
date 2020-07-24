package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"login-provider/internal/config"
	"login-provider/internal/handler"
	"login-provider/internal/hydra"
	"login-provider/internal/middleware"
)

func Serve(cmd *cobra.Command, args []string) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorrelationId())
	router.Use(middleware.RequestId())
	router.Use(middleware.Logger())
	router.LoadHTMLGlob("web/templates/*")

	// TODO: refactor this part and move it closer to the handler functions
	hf, err := hydra.NewClientFactory(config.HydraAdminUrl())
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

	addr := config.Address()

	if tlsConfig, err := config.TlsConfig(); err == nil {
		log.Info().
			Msg("Listening and serving HTTPS on " + addr)
		router.RunTLS(addr, tlsConfig.CertFile, tlsConfig.KeyFile)
	} else {
		log.Info().
			Msg("Listening and serving HTTP on " + addr)
		router.Run(addr)
	}
}