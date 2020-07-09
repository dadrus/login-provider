package hydra

import (
	"github.com/ory/hydra-client-go/client"
	"net/url"
)

type ClientFactory struct {
	adminUrl *url.URL
}

func NewClientFactory(adminUrl string) *ClientFactory {
	adminURL, err := url.Parse(adminUrl)
	if err != nil {
		panic(err)
	}

	return &ClientFactory{
		adminUrl: adminURL,
	}
}

func (h *ClientFactory) NewClient() *client.OryHydra {
	// TODO: configure the trust store to be used. Without the system wide trust store will be used
	// this is not what we want

	// TODO: configure hydra logger to the level this application is configured to use

	return client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes: []string{h.adminUrl.Scheme}, Host: h.adminUrl.Host, BasePath: h.adminUrl.Path})
}
