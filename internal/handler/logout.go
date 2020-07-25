package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/rs/zerolog/log"
	"login-provider/internal/config"
	"login-provider/internal/hydra"
	"net/http"
)

type logoutForm struct {
	Challenge      string `form:"challenge" binding:"required"`
	LogoutApproved bool   `form:"logout_approved"`
}

func ShowLogoutPage(h *hydra.ClientFactory, conf config.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.Ctx(c.Request.Context())

		var logoutChallenge string
		// the challenge is used to fetch information about consent requests in hydra
		if logoutChallenge = c.Query("logout_challenge"); len(logoutChallenge) == 0 {
			HandleBadRequest(c, conf)
			return
		}

		client := h.NewClient(c)
		_, err := client.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().
			WithLogoutChallenge(logoutChallenge))
		if err != nil {
			logger.Err(err).Msg("Error while communicating with hydra to get logout request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
			return
		}

		c.HTML(http.StatusOK, "logout.html", gin.H{
			"title":     "Logout",
			"challenge": logoutChallenge,
		})
	}
}

func Logout(hf *hydra.ClientFactory, _ config.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.Ctx(c.Request.Context())

		var logoutData logoutForm
		if err := c.ShouldBind(&logoutData); err != nil {
			logger.Err(err).Msg("Failed to parse data from submitted logout form")
			c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
			return
		}

		client := hf.NewClient(c)

		if !logoutData.LogoutApproved {
			_, err := client.Admin.RejectLogoutRequest(admin.NewRejectLogoutRequestParams().
				WithLogoutChallenge(logoutData.Challenge))
			if err != nil {
				logger.Err(err).Msg("Error while communicating with hydra to reject logout request")
				// TODO: This is an internal error (hydra not available, the request is malformed, etc)
				// So we have to redirect to "something went wrong page - please contact the admin"
				c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
				return
			}

			// TODO: where to redirect
			c.Redirect(302, "where to redirect???")
			return
		}

		response, err := client.Admin.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParams().
			WithLogoutChallenge(logoutData.Challenge))
		if err != nil {
			logger.Err(err).Msg("Error while communicating with hydra to accept logout request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
			return
		}

		c.Redirect(302, response.Payload.RedirectTo)
	}
}
