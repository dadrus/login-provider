package view

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"login-provider/internal/config"
	"net/http"
)

func HandleBadRequest(c *gin.Context) {
	c.HTML(http.StatusBadRequest,
		"login.html",
		gin.H{"title": "Login", "register_url": viper.GetString(config.RegisterUrl)})
}
