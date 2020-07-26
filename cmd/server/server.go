package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"login-provider/internal/config"
	"login-provider/internal/handler"
	"login-provider/internal/logging"
	"login-provider/internal/middleware"
	"os"
	"strings"
)

func Serve(cmd *cobra.Command, args []string) {
	conf := config.NewConfiguration()
	logging.ConfigureLogging(conf)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CorrelationId())
	router.Use(middleware.RequestId())
	router.Use(middleware.Logger())
	if strings.HasSuffix(os.Getenv("PWD"), "cmd") {
		// because of root_test.go
		router.LoadHTMLGlob("../web/templates/*")
	} else {
		router.LoadHTMLGlob("web/templates/*")
	}

	handler.RegisterRoutes(router, conf)

	addr := conf.Address()
	if tlsConfig, err := conf.TlsConfig(); err == nil {
		log.Info().
			Msg("Listening and serving HTTPS on " + addr)
		router.RunTLS(addr, tlsConfig.CertFile, tlsConfig.KeyFile)
	} else {
		log.Info().
			Msg("Listening and serving HTTP on " + addr)
		router.Run(addr)
	}
}