package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"login-provider/internal/config"
	"login-provider/internal/hydra"
)

func RegisterRoutes(e *gin.Engine, conf config.Configuration) {
	hf, err := hydra.NewClientFactory(conf)
	if err != nil {
		l := log.With().Err(err).Logger()
		l.Fatal().Msg("Failed to create hydra client factory")
	}

	e.GET("/login", ShowLoginPage(hf, conf))
	e.POST("/login", Login(hf, conf))
	e.GET("/consent", ShowConsentPage(hf, conf))
	e.POST("/consent", Consent(hf, conf))
	e.GET("/logout", ShowLogoutPage(hf, conf))
	e.POST("/logout", Logout(hf, conf))

	e.GET("/health/alive", Alive)
	e.GET("/health/ready", Ready(hf))
}
