package logger

import (
	"github.com/collinzh/chatbot-messenger/conf"
	"github.com/rs/zerolog"
)

var (
	logger zerolog.Logger
)

func Initialize() {
	if conf.Debug() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	out := zerolog.NewConsoleWriter()
	logger = zerolog.New(out)
	logger = logger.With().Caller().Logger()
}

func Logger() *zerolog.Logger {
	return &logger
}
