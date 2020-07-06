package view

import (
	"github.com/ory/hydra-client-go/client"
	"net/url"
)

type HydraClientFactory struct {
	adminUrl *url.URL
}

func NewHydraClientFactory(adminUrl string) *HydraClientFactory {
	adminURL, err := url.Parse(adminUrl)
	if err != nil {
		panic(err)
	}

	return &HydraClientFactory{
		adminUrl: adminURL,
	}
}

func (h *HydraClientFactory) newClient() *client.OryHydra {
	// TODO: configure the trust store to be used. Without the system wide trust store will be used
	// this is not what we want

	// TODO: configure hydra logger to the level this application is configured to use

	return client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes: []string{h.adminUrl.Scheme}, Host: h.adminUrl.Host, BasePath: h.adminUrl.Path})
}
