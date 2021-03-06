package hydra

import (
	"context"
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ory/hydra-client-go/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"login-provider/internal/config"
	"net/url"
)

type ClientFactory struct {
	transport *httptransport.Runtime
}

func NewClientFactory(conf config.Configuration) (*ClientFactory, error) {
	url, err := url.Parse(conf.HydraAdminUrl())
	if err != nil {
		return nil, err
	}

	factory := &ClientFactory{}

	if caFile, err := conf.TlsTrustStore(); err != nil {
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

	factory.transport.SetDebug(conf.LogLevel() == zerolog.DebugLevel)

	return factory, nil
}

func (cf *ClientFactory) NewClient(ctx context.Context) *client.OryHydra {
	logger := log.Ctx(ctx)
	cf.transport.SetLogger(zeroLogLogger{logger})
	return client.New(cf.transport, nil)
}

type zeroLogLogger struct{
	logger *zerolog.Logger
}

func (l zeroLogLogger) Printf(format string, args ...interface{}) {
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}

	l.logger.Info().Msg(fmt.Sprintf(format, args))
}

func (l zeroLogLogger) Debugf(format string, args ...interface{}) {
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}

	l.logger.Debug().Msg(fmt.Sprintf(format, args))
}
