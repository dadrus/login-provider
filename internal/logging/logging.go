package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"login-provider/internal/config"
	"os"
	"strings"
)

func ConfigureLogging(conf config.Configuration) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "_level_name"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}
	zerolog.MessageFieldName = "short_message"
	zerolog.ErrorFieldName = "full_message"
	zerolog.CallerFieldName = "_caller"
	zerolog.SetGlobalLevel(conf.LogLevel())

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
		Logger().Hook(zerolog.HookFunc(
			func(e *zerolog.Event, level zerolog.Level, message string) {
				if level != zerolog.NoLevel {
					e.Int("level", toSyslogLevel(level))
				}
			}))
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
