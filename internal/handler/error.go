package handler

import (
	"github.com/gin-gonic/gin"
	"login-provider/internal/config"
	"net/http"
)

func HandleBadRequest(c *gin.Context, conf config.Configuration) {
	c.HTML(http.StatusBadRequest,
		"login.html",
		gin.H{"title": "Login", "register_url": conf.RegisterUrl()})
}
