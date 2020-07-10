package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"login-provider/internal/hydra"
	"net/http"
)

type logoutForm struct {
	Challenge      string `form:"challenge" binding:"required"`
	LogoutApproved bool   `form:"logout_approved"`
}

func ShowLogoutPage(h *hydra.ClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logoutChallenge string
		// the challenge is used to fetch information about consent requests in hydra
		if logoutChallenge = c.Query("logout_challenge"); len(logoutChallenge) == 0 {
			HandleBadRequest(c)
			return
		}

		client := h.NewClient(c)
		_, err := client.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().
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

func Logout(hf *hydra.ClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logoutData logoutForm
		if err := c.ShouldBind(&logoutData); err != nil {
			c.HTML(http.StatusBadRequest, "logout.html", gin.H{"title": "Logout"})
			return
		}

		client := hf.NewClient(c)

		if !logoutData.LogoutApproved {
			_, err := client.Admin.RejectLogoutRequest(admin.NewRejectLogoutRequestParams().
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

		response, err := client.Admin.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParams().
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
