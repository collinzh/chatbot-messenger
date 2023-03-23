package storage

import (
	"context"
	"github.com/collinzh/chatbot-messenger/conf"
	"github.com/collinzh/chatbot-messenger/types"
)

type Storage interface {
	Connect() error
	Shutdown()

	FindChat(authCode string, ctx context.Context) (*types.Chat, error)
	SaveChat(authCode string, messages []types.Message, ctx context.Context) error
	Prune() error
}

func New() Storage {
	if conf.MongoURI() == "" {
		return newMemoryStorage()
	} else {
		return newMongoStorage()
	}
}
