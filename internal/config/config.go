package config

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	// TODO: use home documents to retrieve the following three values using the forth one
	registerUrl     = "register_url"
	authenticateUrl = "authenticate_url"
	hydraAdminUrl   = "hydra_admin_url"
	rootHomeUrl = "root_home_url"

	tlsKeyFile = "tls.key"
	tlsCertFile = "tls.cert"
	tlsTrustStoreFile = "tls.trust_store"

	logLevel = "log.level"

	host = "host"
	port = "port"
)

type Configuration interface {
	// TODO: update methods returning Urls to return URL type and error
	Address() string
	TlsConfig() (*TlsConfig, error)
	TlsTrustStore() (string, error)
	RegisterUrl() string
	AuthenticateUrl() string
	HydraAdminUrl() string
	LogLevel() zerolog.Level
}

type TlsConfig struct {
	KeyFile  string
	CertFile string
}

// Loads and reads the config and environment variables if set
func Load(file *string) func() {
	return func() {
		if *file != "" {
			// enable ability to specify config file via flag
			viper.SetConfigFile(*file)
		} else {
			viper.SetConfigType("yaml")
			viper.SetConfigName("config")
			viper.AddConfigPath(".")
		}

		viper.SetDefault(logLevel, "info")
		viper.SetDefault(port, "8080")
		viper.SetDefault(host, "127.0.0.1")

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("No config file found.")
			os.Exit(-1)
		}
	}
}

type configuration struct {}

func NewConfiguration() Configuration {
	return &configuration{}
}

func (c *configuration) Address() string {
	return viper.GetString(host) + ":" + viper.GetString(port)
}

func (c *configuration) TlsConfig() (*TlsConfig, error) {
	tlsKeyFile := viper.GetString(tlsKeyFile)
	if len(tlsKeyFile) == 0 {
		return nil, errors.New("no TLS key configured")
	}
	if _, err := os.Stat(tlsKeyFile); err != nil {
		return nil, errors.New("configured TLS key not available")
	}

	tlsCertFile := viper.GetString(tlsCertFile)
	if len(tlsCertFile) == 0 {
		return nil, errors.New("no TLS cert configured")
	}
	if _, err := os.Stat(tlsCertFile); err != nil {
		return nil, errors.New("configured TLS cert not available")
	}

	return &TlsConfig{
		KeyFile:  tlsKeyFile,
		CertFile: tlsCertFile,
	}, nil
}

func (c *configuration) TlsTrustStore() (string, error) {
	value := viper.GetString(tlsTrustStoreFile)
	if len(value) == 0 {
		return "", errors.New("no TLS key configured")
	}
	if _, err := os.Stat(value); err != nil {
		return "", errors.New("configured TLS key not available")
	}
	return value, nil
}

func (c *configuration) RegisterUrl() string  {
	return viper.GetString(registerUrl)
}

func (c *configuration) HydraAdminUrl() string  {
	return viper.GetString(hydraAdminUrl)
}

func (c *configuration) AuthenticateUrl() string  {
	return viper.GetString(authenticateUrl)
}

func (c *configuration) LogLevel() zerolog.Level  {
	switch viper.GetString(logLevel) {
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}


