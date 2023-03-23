package main

import (
	"github.com/collinzh/chatbot-messenger/bot"
	"github.com/collinzh/chatbot-messenger/conf"
	"github.com/collinzh/chatbot-messenger/logger"
	"github.com/collinzh/chatbot-messenger/server"
	"github.com/collinzh/chatbot-messenger/storage"
	"net/http"
)

func main() {
	conf.ParseConfig()
	logger.Initialize()

	s := storage.New()
	if err := s.Connect(); err != nil {
		logger.Logger().Fatal().Err(err).Msg("Cannot connect to storage")
	}

	ai := bot.New(s)

	err := server.New(ai).RunAndBlock()
	if err != http.ErrServerClosed {
		logger.Logger().Fatal().Err(err).Msg("Error creating web server")
	}

	logger.Logger().Info().Msg("Shutting down")
}
