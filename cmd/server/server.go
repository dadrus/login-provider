package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"login-provider/internal/config"
	"login-provider/internal/handler"
	"login-provider/internal/middleware"
)

func Serve(cmd *cobra.Command, args []string) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorrelationId())
	router.Use(middleware.RequestId())
	router.Use(middleware.Logger())
	router.LoadHTMLGlob("web/templates/*")

	handler.RegisterRoutes(router)

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