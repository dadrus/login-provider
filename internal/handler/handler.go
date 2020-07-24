package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"login-provider/internal/config"
	"login-provider/internal/hydra"
)

func RegisterRoutes(e *gin.Engine) {
	hf, err := hydra.NewClientFactory(config.HydraAdminUrl())
	if err != nil {
		l := log.With().Err(err).Logger()
		l.Fatal().Msg("Failed to create hydra client factory")
	}

	e.GET("/login", ShowLoginPage(hf))
	e.POST("/login", Login(hf))
	e.GET("/consent", ShowConsentPage(hf))
	e.POST("/consent", Consent(hf))
	e.GET("/logout", ShowLogoutPage(hf))
	e.POST("/logout", Logout(hf))
}
