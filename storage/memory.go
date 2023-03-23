package storage

import (
	"context"
	"github.com/collinzh/chatbot-messenger/types"
)

type memoryStorageImpl struct {
	persist map[string]string
}

func (m *memoryStorageImpl) FindChat(authCode string, ctx context.Context) (*types.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (m *memoryStorageImpl) SaveChat(authCode string, messages []types.Message, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *memoryStorageImpl) Prune() error {
	//TODO implement me
	panic("implement me")
}

func newMemoryStorage() Storage {
	return &memoryStorageImpl{}
}

func (m *memoryStorageImpl) Connect() error {
	return nil
}

func (m *memoryStorageImpl) Shutdown() {

}
