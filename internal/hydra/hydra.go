package hydra

import (
	"fmt"
	"github.com/gin-gonic/gin"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ory/hydra-client-go/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"login-provider/internal/config"
	"net/url"
)

type ClientFactory struct {
	transport *httptransport.Runtime
}

func NewClientFactory(adminUrl string) (*ClientFactory, error) {
	url, err := url.Parse(adminUrl)
	if err != nil {
		return nil, err
	}

	factory := &ClientFactory{}

	caFile := viper.GetString(config.TlsTrustStoreFile)
	if caFile == "" {
		log.Info().Msg("No explicit trust store configured. Falling back to a system-wide one")
		// if a specific trust store is not specified, we'll rely on the system-wide trust store
		factory.transport = httptransport.New(url.Host, url.Path, []string{url.Scheme})
	} else {
		log.Info().Msg("Explicit trust store configured. Using it")
		// if a specific trust store has been specified use it instead fo the the system wide one
		tlsClient, err := httptransport.TLSClient(httptransport.TLSClientOptions{
			CA: caFile,
		})

		if err != nil {
			return nil, err
		}

		factory.transport = httptransport.NewWithClient(url.Host, url.Path, []string{url.Scheme}, tlsClient)
	}

	factory.transport.SetDebug(viper.GetString(config.LogLevel) == "debug")

	return factory, nil
}

func (cf *ClientFactory) NewClient(c *gin.Context) *client.OryHydra {
	logger := log.Ctx(c.Request.Context())
	cf.transport.SetLogger(ZeroLogLogger{logger})
	return client.New(cf.transport, nil)
}

type ZeroLogLogger struct{
	logger *zerolog.Logger
}

func (l ZeroLogLogger) Printf(format string, args ...interface{}) {
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}

	l.logger.Info().Msg(fmt.Sprintf(format, args))
}

func (l ZeroLogLogger) Debugf(format string, args ...interface{}) {
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}

	l.logger.Debug().Msg(fmt.Sprintf(format, args))
}
