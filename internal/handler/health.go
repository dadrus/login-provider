package handler

import (
	"github.com/gin-gonic/gin"
	"login-provider/internal/hydra"
	"net/http"
)

type healthStatus struct {
	// Status always contains "ok".
	Status string `json:"status"`
}

type readyStatus struct {
	healthStatus
	Errors map[string]string `json:"errors,omitempty""`
}

func Alive(c *gin.Context) {
	c.JSON(http.StatusOK, &healthStatus{Status: "Ok"})
}

func Ready(hf *hydra.ClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement readyness probe. use hydras /health/ready end point for that
	}
}
