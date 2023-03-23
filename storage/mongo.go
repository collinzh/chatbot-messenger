package storage

import (
	"context"
	"errors"
	"github.com/collinzh/chatbot-messenger/conf"
	"github.com/collinzh/chatbot-messenger/logger"
	"github.com/collinzh/chatbot-messenger/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"time"
)

type mongodbStorageImpl struct {
	mongoClient *mongo.Client
	database    *mongo.Database
	dbName      string
}

const (
	collectionChat = "chats"
)

var (
	ErrNotFound = errors.New("no such auth code")
	ErrInternal = errors.New("unable to decode document")
)

func (m *mongodbStorageImpl) FindChat(authCode string, ctx context.Context) (*types.Chat, error) {
	collection := m.database.Collection(collectionChat)

	result := &types.Chat{}
	err := collection.FindOne(ctx, bson.D{{"auth_code", authCode}}).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Logger().Debug().Err(err).Str("auth_code", authCode).Msg("Cannot find chats")
			return nil, ErrNotFound
		}
		logger.Logger().Error().Err(err).Str("auth_code", authCode).Msgf("Unable to lookup chat by auth code %s", authCode)
		return nil, ErrInternal
	}

	return result, nil
}

func (m *mongodbStorageImpl) SaveChat(authCode string, messages []types.Message, ctx context.Context) error {
	collection := m.database.Collection(collectionChat)

	update := bson.D{{"$set", bson.D{{"last_updated", time.Now().UnixMilli()}, {"messages", messages}}}}

	res, err := collection.UpdateOne(ctx, bson.D{{"auth_code", authCode}}, update)

	if err != nil {
		logger.Logger().Error().Err(err).Str("auth_code", authCode).Msg("Error saving chats")
		return ErrInternal
	}

	if res.ModifiedCount != 1 {
		logger.Logger().Warn().Str("auth_code", authCode).Msg("Chat not updated")
	}

	return nil
}

func (m *mongodbStorageImpl) Prune() error {
	// for chats that have not been updated in a while
	deleteUntil := time.Now().UnixMilli() - conf.Retention().Milliseconds()
	filter := bson.D{{"last_updated", bson.D{{"$lt", deleteUntil}}}}

	// Remove message history
	update := bson.D{{"$set", bson.D{{"messages", []string{}}}}}

	_, err := m.database.Collection(collectionChat).UpdateMany(context.Background(), filter, update)
	if err != nil {
		logger.Logger().Err(err).Msg("Unable to prune storage")
	}

	return nil
}

func newMongoStorage() Storage {
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.MongoURI()))
	if err != nil {
		panic(err)
	}

	// MongoDB client has already parsed and validated the connection string. No need to check for error here
	cs, err := connstring.Parse(conf.MongoURI())

	return &mongodbStorageImpl{mongoClient: client, dbName: cs.Database}
}

func (m *mongodbStorageImpl) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.mongoClient.Connect(ctx)
	if err != nil {
		return err
	}
	m.database = m.mongoClient.Database(m.dbName)

	return nil
}

func (m *mongodbStorageImpl) Shutdown() {
	_ = m.mongoClient.Disconnect(context.Background())
}
