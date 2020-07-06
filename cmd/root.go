package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"login-provider/cmd/server"
	"login-provider/internal/config"
	"os"
	"strings"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "login-provider",
	Short: "Hydra login provider",
	Long:  "Hydra login provider offering UI controls for OIDC login, consent and logout flows",
	Run:   server.Serve,
}

func init() {
	cobra.OnInitialize(setConfigDefaults, readConfig, configureLogging)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $PWD/config.yaml)")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func setConfigDefaults() {
	viper.SetDefault(config.LogLevel, "info")
	viper.SetDefault(config.Port, "8080")
	viper.SetDefault(config.CookieSameSiteMode, "Lax")
}

// Reads config and environment variables if set
func readConfig() {
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		/*if _, err := os.Stat("./config.yml"); err != nil {
			_, _ = os.Create("./config.yml")
		}*/

		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found.")
		os.Exit(-1)
	}
}

func configureLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch viper.GetString("log.level") {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
