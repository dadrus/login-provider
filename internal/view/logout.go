package view

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
)

type logoutForm struct {
	Challenge      string `form:"challenge" binding:"required"`
	LogoutApproved bool   `form:"logout_approved"`
}

func ShowLogoutPage(h *HydraClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logoutChallenge string
		// the challenge is used to fetch information about consent requests in hydra
		if logoutChallenge = c.Query("logout_challenge"); len(logoutChallenge) == 0 {
			HandleBadRequest(c)
			return
		}

		hydra := h.newClient()
		_, err := hydra.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().
			WithLogoutChallenge(logoutChallenge))
		if err != nil {
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

func Logout(hf *HydraClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logoutData logoutForm
		if err := c.ShouldBind(&logoutData); err != nil {
			c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
			return
		}

		hydra := hf.newClient()

		if !logoutData.LogoutApproved {
			_, err := hydra.Admin.RejectLogoutRequest(admin.NewRejectLogoutRequestParams().
				WithLogoutChallenge(logoutData.Challenge))
			if err != nil {
				// TODO: This is an internal error (hydra not available, the request is malformed, etc)
				// So we have to redirect to "something went wrong page - please contact the admin"
				c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
				return
			}

			// TODO: where to redirect
			c.Redirect(302, "where to redirect???")
			return
		}

		response, err := hydra.Admin.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParams().
			WithLogoutChallenge(logoutData.Challenge))
		if err != nil {
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
			return
		}

		c.Redirect(302, response.Payload.RedirectTo)
	}
}
