package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/rs/zerolog/log"
	"login-provider/internal/client_meta"
	"login-provider/internal/hydra"
	"login-provider/internal/profile_api"
	"net/http"
	"time"
)

type consentForm struct {
	Challenge       string   `form:"challenge" binding:"required"`
	GrantedScopes   []string `form:"granted_scopes[]"`
	Remember        bool     `form:"remember"`
	ConsentApproved bool     `form:"consent_approved"`
}

func ShowConsentPage(hf *hydra.ClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var consentChallenge string

		// the challenge is used to fetch information about consent requests in hydra
		if consentChallenge = c.Query("consent_challenge"); len(consentChallenge) == 0 {
			log.Warn().Msg("No consent challenge provided")
			HandleBadRequest(c)
			return
		}

		client := hf.NewClient()

		response, err := client.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().
			WithConsentChallenge(consentChallenge))
		if err != nil {
			log.Err(err).Msg("Error while communicating with hydra to get consent request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest, "consent.html", gin.H{"title": "Consent"})
			return
		}

		// raw data contains the url to the users profile, as well as all the data, which can be retrieved
		// from that endpont. So parse it and set the values in AccessToken and IDToken accordingly taking
		// granted scopes into account
		authResponse := &profile_api.AuthenticationResponse{}
		err = authResponse.Unmarshal(response.Payload.Context)

		// check whether consent is required for given client and which consent values are required
		info := &client_meta.ClientMetaInfo{}
		err = info.Unmarshal(response.Payload.Client.Metadata)

		if response.Payload.Skip || !info.AskConsent {
			// grant login request
			response, err := client.Admin.AcceptConsentRequest(
				admin.NewAcceptConsentRequestParams().
					WithConsentChallenge(consentChallenge).
					WithBody(&models.AcceptConsentRequest{
						GrantAccessTokenAudience: response.Payload.RequestedAccessTokenAudience,
						GrantScope:               response.Payload.RequestedScope,
						RememberFor:              3600,
						HandledAt:                models.NullTime(time.Now()),
						Session: &models.ConsentRequestSession{
							IDToken: authResponse.CreateIdTokenClaims(response.Payload.RequestedScope),
						},
					}))

			if err != nil {
				log.Err(err).Msg("Error while communicating with hydra to accept consent request")
				// TODO: This is an internal error (hydra not available, the request is malformed, etc)
				// So we have to redirect to "something went wrong page - please contact the admin"
				c.HTML(http.StatusBadRequest, "consent.html", gin.H{"title": "Consent"})
				return
			}

			c.Redirect(302, response.Payload.RedirectTo)
			return
		}

		// look which requested scopes are mandatory
		scopeInfos := info.CreateScopeInfos(response.Payload.RequestedScope)

		// If we are here render Consent page
		c.HTML(http.StatusOK, "consent.html", gin.H{
			"title":           "Consent",
			"challenge":       consentChallenge,
			"requestedScopes": scopeInfos,
			"user":            authResponse.User.UserName,
			"client":          response.Payload.Client,
		})
	}
}

func Consent(hf *hydra.ClientFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var consentData consentForm
		if err := c.ShouldBind(&consentData); err != nil {
			log.Err(err).Msg("Failed to parse data from submitted consent form")
			c.HTML(http.StatusBadRequest, "consent.html", gin.H{"title": "Consent"})
			return
		}

		client := hf.NewClient()

		if !consentData.ConsentApproved {
			response, err := client.Admin.RejectConsentRequest(admin.NewRejectConsentRequestParams().
				WithConsentChallenge(consentData.Challenge).
				WithBody(&models.RejectRequest{
					Error:     "User rejected consent",
					ErrorHint: "consent_rejected",
				}))
			if err != nil {
				log.Err(err).Msg("Error while communicating with hydra to reject consent request")
				// TODO: This is an internal error (hydra not available, the request is malformed, etc)
				// So we have to redirect to "something went wrong page - please contact the admin"
				c.HTML(http.StatusBadRequest, "consent.html", gin.H{"title": "Consent"})
				return
			}

			c.Redirect(302, response.Payload.RedirectTo)
			return
		}

		gcr, err := client.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().
			WithConsentChallenge(consentData.Challenge))
		if err != nil {
			log.Err(err).Msg("Error while communicating with hydra to get consent request")
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest, "consent.html", gin.H{"title": "Consent"})
			return
		}

		ar := &profile_api.AuthenticationResponse{}
		err = ar.Unmarshal(gcr.Payload.Context)

		cmi := &client_meta.ClientMetaInfo{}
		err = cmi.Unmarshal(gcr.Payload.Client.Metadata)

		grantedScopes := append(consentData.GrantedScopes, cmi.MandatoryScopes...)

		acr, err := client.Admin.AcceptConsentRequest(
			admin.NewAcceptConsentRequestParams().
				WithConsentChallenge(consentData.Challenge).
				WithBody(&models.AcceptConsentRequest{
					GrantAccessTokenAudience: gcr.Payload.RequestedAccessTokenAudience,
					GrantScope:               append(consentData.GrantedScopes, cmi.MandatoryScopes...),
					RememberFor:              3600,
					Remember:                 consentData.Remember,
					HandledAt:                models.NullTime(time.Now()),
					Session: &models.ConsentRequestSession{
						IDToken: ar.CreateIdTokenClaims(grantedScopes),
					},
				}))
		if err != nil {
			// TODO: This is an internal error (hydra not available, the request is malformed, etc)
			// So we have to redirect to "something went wrong page - please contact the admin"
			c.HTML(http.StatusBadRequest, "consent.html", gin.H{"title": "Consent"})
			return
		}
		c.Redirect(302, acr.Payload.RedirectTo)
	}
}
