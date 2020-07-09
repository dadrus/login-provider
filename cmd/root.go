package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "_level_name"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}
	zerolog.MessageFieldName = "short_message"
	zerolog.ErrorFieldName = "full_message"
	zerolog.CallerFieldName = "_caller"

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

	log.Logger = zerolog.New(os.Stdout).With().
		Str("version", "1.1").
		Timestamp().Caller().
		Logger().Hook(zerolog.HookFunc(LevelHook))
}

func LevelHook(e *zerolog.Event, level zerolog.Level, message string) {
	if level != zerolog.NoLevel {
		e.Int("level", toSyslogLevel(level))
	}
}

func toSyslogLevel(level zerolog.Level) int {
	switch level {
	case zerolog.DebugLevel, zerolog.TraceLevel:
		return 7
	case zerolog.InfoLevel:
		return 6
	case zerolog.WarnLevel:
		return 4
	case zerolog.ErrorLevel:
		return 3
	case zerolog.FatalLevel:
		return 2
	case zerolog.PanicLevel:
		return 1
	default:
		return 0
	}
}
