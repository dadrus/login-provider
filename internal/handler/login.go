package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/rs/zerolog/log"
	"login-provider/internal/config"
	"login-provider/internal/hydra"
	"login-provider/internal/profile_api"
	"net/http"
	"net/url"
	"strconv"
)

// TODO annotate the handlers for generating OpenAPI spec out of it, e.g. by using https://github.com/go-swagger/go-swagger

type loginForm struct {
	Challenge string `form:"challenge" binding:"required"`
	Email     string `form:"email" binding:"required"`
	Password  string `form:"password" binding:"required"`
	Remember  bool   `form:"remember"`
}

func ShowLoginPage(hf *hydra.ClientFactory, conf config.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.Ctx(c.Request.Context())

		logger.Debug().Msg("Showing login page")
		var loginChallenge string

		// the challenge is used to fetch information about login requests in hydra
		if loginChallenge = c.Query("login_challenge"); len(loginChallenge) == 0 {
			logger.Warn().Msg("No login challenge provided")
			HandleBadRequest(c, conf)
			return
		}

		errorMessage := c.Query("error")

		client := hf.NewClient(c.Request.Context())
		// get info about the login request for the given challenge
		response, err := client.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().
			WithLoginChallenge(loginChallenge))
		if err != nil {
			logger.Err(err).Msg("Error while communicating with hydra to get new login request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			HandleBadRequest(c, conf)
			return
		}

		// if hydra was already able to authenticate the user, Skip will be true
		// and we don't need to authenticate the user again
		if response.Payload.Skip {
			logger.Debug().Msg("User authentication skipped")

			// grant login request
			response, err := client.Admin.AcceptLoginRequest(
				admin.NewAcceptLoginRequestParams().
					WithLoginChallenge(loginChallenge).
					WithBody(&models.AcceptLoginRequest{Subject: &response.Payload.Subject}))
			if err != nil {
				logger.Err(err).Msg("Error while communicating with hydra to accept login request")
				// TODO: This is an internal error (hydra not available, the request is malformed, etc)
				// So we have to redirect to "something went wrong page - please contact the admin"
				HandleBadRequest(c, conf)
				return
			}

			c.Redirect(302, response.Payload.RedirectTo)
			return
		}

		// If we are here render Login page
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":        "Login",
			"challenge":    loginChallenge,
			"register_url": conf.RegisterUrl(),
			"error":        errorMessage,
		})
	}
}

func Login(hf *hydra.ClientFactory, conf config.Configuration) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.Ctx(c.Request.Context())

		var loginData loginForm
		if err := c.ShouldBind(&loginData); err != nil {
			logger.Err(err).Msg("Failed to parse data from submitted login form")
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"title": "Login", "register_url": conf.RegisterUrl()})
			return
		}

		client := hf.NewClient(c.Request.Context())
		authResponse, err := profile_api.AuthenticateUser(conf.AuthenticateUrl(), loginData.Email, loginData.Password)
		if err != nil {
			l := logger.With().Err(err).Logger()
			l.Warn().Msg("User authentication failed")
			params := url.Values{}
			params.Add("login_challenge", loginData.Challenge)
			params.Add("error", "Invalid user name or password")
			c.Redirect(302, "/login?" + params.Encode())
			return
		}

		subjectId := strconv.Itoa(authResponse.User.ID)

		// login successful
		response, err := client.Admin.AcceptLoginRequest(admin.NewAcceptLoginRequestParams().
			WithLoginChallenge(loginData.Challenge).
			WithBody(&models.AcceptLoginRequest{
				Acr:         "0",
				Context:     authResponse,
				Remember:    loginData.Remember,
				RememberFor: 3600,
				Subject:     &subjectId,
			}))
		if err != nil {
			logger.Err(err).Msg("Error while communicating with hydra to accept login request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest,
				"login.html",
				gin.H{"title": "Login", "error": "Login failed", "register_url": conf.RegisterUrl()})
			return
		}

		c.Redirect(302, response.Payload.RedirectTo)
	}
}