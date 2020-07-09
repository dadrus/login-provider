package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"login-provider/internal/config"
	"login-provider/internal/profile_api"
	"net/http"
	"net/url"
	"strconv"
)

type loginForm struct {
	Challenge string `form:"challenge" binding:"required"`
	Email     string `form:"email" binding:"required"`
	Password  string `form:"password" binding:"required"`
	Remember  bool   `form:"remember"`
}

func ShowLoginPage(hf *HydraClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(zerolog.Logger)

		log.Debug().Msg("Showing login page")
		var loginChallenge string

		// the challenge is used to fetch information about login requests in hydra
		if loginChallenge = c.Query("login_challenge"); len(loginChallenge) == 0 {
			log.Warn().Msg("No login challenge provided")
			HandleBadRequest(c)
			return
		}

		errorMessage := c.Query("error")

		hydra := hf.newClient()
		// get info about the login request for the given challenge
		response, err := hydra.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().
			WithLoginChallenge(loginChallenge))
		if err != nil {
			log.Err(err).Msg("Error while communicating with hydra to get new login request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			HandleBadRequest(c)
			return
		}

		// if hydra was already able to authenticate the user, Skip will be true
		// and we don't need to authenticate the user again
		if response.Payload.Skip {
			log.Debug().Msg("User authentication skipped")

			// grant login request
			response, err := hydra.Admin.AcceptLoginRequest(
				admin.NewAcceptLoginRequestParams().
					WithLoginChallenge(loginChallenge).
					WithBody(&models.AcceptLoginRequest{Subject: &response.Payload.Subject}))
			if err != nil {
				log.Err(err).Msg("Error while communicating with hydra to accept login request")
				// TODO: This is an internal error (hydra not available, the request is malformed, etc)
				// So we have to redirect to "something went wrong page - please contact the admin"
				HandleBadRequest(c)
				return
			}

			c.Redirect(302, response.Payload.RedirectTo)
			return
		}

		// If we are here render Login page
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":        "Login",
			"challenge":    loginChallenge,
			"register_url": viper.GetString(config.RegisterUrl),
			"error":        errorMessage,
		})
	}
}

func Login(hf *HydraClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(zerolog.Logger)

		var loginData loginForm
		if err := c.ShouldBind(&loginData); err != nil {
			log.Err(err).Msg("Failed to parse data from submitted login form")
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"title": "Login", "register_url": viper.GetString(config.RegisterUrl)})
			return
		}

		hydra := hf.newClient()
		authResponse, err := profile_api.AuthenticateUser(viper.GetString(config.AuthenticateUrl), loginData.Email, loginData.Password)
		if err != nil {
			params := url.Values{}
			params.Add("login_challenge", loginData.Challenge)
			params.Add("error", "Invalid user name or password")
			c.Redirect(302, "/login?" + params.Encode())
			return
		}

		subjectId := strconv.Itoa(authResponse.User.ID)

		// login successful
		response, err := hydra.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(loginData.Challenge).
			WithBody(&models.AcceptLoginRequest{
				Acr:         "0",
				Context:     authResponse,
				Remember:    loginData.Remember,
				RememberFor: 3600,
				Subject:     &subjectId,
			}))
		if err != nil {
			log.Err(err).Msg("Error while communicating with hydra to accept login request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest,
				"login.html",
				gin.H{"title": "Login", "error": "Login failed", "register_url": viper.GetString(config.RegisterUrl)})
			return
		}

		c.Redirect(302, response.Payload.RedirectTo)
	}
}