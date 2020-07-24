package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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
	cobra.OnInitialize(config.Load(cfgFile), configureLogging)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $PWD/config.yaml)")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
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
	zerolog.SetGlobalLevel(config.LogLevel())

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
		fmt.Println("Failed to retrieve the hostname: " + err.Error())
	}

	log.Logger = zerolog.New(os.Stdout).With().
		Str("version", "1.1").
		Str("host", hostname).
		Timestamp().
		Caller().
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
