package bot

import (
	"context"
	"errors"
	"github.com/collinzh/chatbot-messenger/conf"
	"github.com/collinzh/chatbot-messenger/logger"
	"github.com/collinzh/chatbot-messenger/storage"
	"github.com/collinzh/chatbot-messenger/types"
	"github.com/sashabaranov/go-openai"
)

type OpenAI interface {
	SubmitChat(msg, authCode string, ctx context.Context) (string, error)
}

type openAiImpl struct {
	client *openai.Client
	sto    storage.Storage
}

var (
	ErrCannotSubmitChat = errors.New("cannot submit chat")
)

func New(sto storage.Storage) OpenAI {
	c := &openAiImpl{client: openai.NewClient(conf.OpenAPIKey()), sto: sto}
	return c
}

// toOpenAiMessages converts a list of messages to OpenAI request format
func toOpenAiMessages(messages []types.Message) []openai.ChatCompletionMessage {
	res := make([]openai.ChatCompletionMessage, len(messages))
	for idx, msg := range messages {
		var role string
		if msg.MsgType == types.BotMessage {
			role = openai.ChatMessageRoleAssistant
		} else {
			role = openai.ChatMessageRoleUser
		}

		res[idx] = openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		}
	}

	return res
}

func (o *openAiImpl) SubmitChat(msg, authCode string, ctx context.Context) (string, error) {
	chat, err := o.sto.FindChat(authCode, ctx)

	if err != nil {
		if err == storage.ErrNotFound {
			return "", types.ErrAuthCode
		}
		logger.Logger().Err(err).Msg("Cannot find chat from data storage")
		return "", ErrCannotSubmitChat
	}

	// Convert existing chats
	messages := toOpenAiMessages(chat.Messages)

	// Append new message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})

	// build and send request
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	}
	res, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		logger.Logger().Err(err).Str("auth_code", authCode).Msg("OpenAPI error")
		return "", ErrCannotSubmitChat
	}

	newChats := append(chat.Messages, types.Message{MsgType: types.UserMessage, Content: msg})
	lastMessage := ""

	for _, reply := range res.Choices {
		newChats = append(newChats, types.Message{
			MsgType: types.BotMessage,
			Content: reply.Message.Content,
		})
		lastMessage = reply.Message.Content
	}

	logger.Logger().Debug().Str("auth_code", authCode).Msgf("ChatGPT replied: %s", res.Object)

	err = o.sto.SaveChat(authCode, newChats, ctx)
	if err != nil {
		logger.Logger().Err(err).Msg("Error saving chats")
		return "", ErrCannotSubmitChat
	}

	return lastMessage, nil
}
