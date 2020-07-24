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
	registerUrl     = "register_url"
	authenticateUrl = "authenticate_url"
	hydraAdminUrl   = "hydra_admin_url"

	tlsKeyFile = "tls.key"
	tlsCertFile = "tls.cert"
	tlsTrustStoreFile = "tls.trust_store"

	logLevel = "log.level"

	host = "host"
	port = "port"
)

// Loads and reads the config and environment variables if set
func Load(file string) func() {
	return func() {
		if file != "" {
			// enable ability to specify config file via flag
			viper.SetConfigFile(file)
		} else {
			viper.SetConfigType("yaml")
			viper.SetConfigName("config")
			viper.AddConfigPath(".")
		}

		viper.SetDefault(logLevel, "info")
		viper.SetDefault(port, "8080")

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("No config file found.")
			os.Exit(-1)
		}
	}
}

func Address() string {
	return viper.GetString(host) + ":" + viper.GetString(port)
}

type tlsConfig struct {
	KeyFile  string
	CertFile string
	TrustStoreFile string
}

func TlsConfig() (*tlsConfig, error) {
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

	tlsTrustStoreFile := viper.GetString(tlsTrustStoreFile)
	if len(tlsTrustStoreFile) != 0 {
		if _, err := os.Stat(tlsTrustStoreFile); err != nil {
			return nil, errors.New("configured TLS trust store not available")
		}
	}


	return &tlsConfig{
		KeyFile:  tlsKeyFile,
		CertFile: tlsCertFile,
		TrustStoreFile: tlsTrustStoreFile,
	}, nil
}

func RegisterUrl() string  {
	return viper.GetString(registerUrl)
}

func HydraAdminUrl() string  {
	return viper.GetString(hydraAdminUrl)
}

func AuthenticateUrl() string  {
	return viper.GetString(authenticateUrl)
}

func LogLevel() zerolog.Level  {
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


