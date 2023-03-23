package server

import (
	"github.com/collinzh/chatbot-messenger/bot"
	"github.com/collinzh/chatbot-messenger/conf"
	"github.com/collinzh/chatbot-messenger/logger"
	"github.com/collinzh/chatbot-messenger/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

type Server interface {
	RunAndBlock() error
}

type serverImpl struct {
	engine *gin.Engine
	openAi bot.OpenAI
}

func New(ai bot.OpenAI) Server {
	s := &serverImpl{openAi: ai}
	err := s.init()
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("Unable to start web server")
		return nil
	}

	return s
}

func (s *serverImpl) init() error {
	if conf.Debug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine = gin.Default()
	s.engine.Use(authorizeRequest())
	s.engine.POST("/api/chat", s.handleChatRequest)

	if len(conf.TrustedProxies()) > 0 {
		if err := s.engine.SetTrustedProxies(conf.TrustedProxies()); err != nil {
			logger.Logger().Fatal().Err(err).Msg("Error configuring trusted proxies")
		}

	}

	return nil
}

func (s *serverImpl) handleChatRequest(c *gin.Context) {
	body := &ChatRequest{}
	if err := c.ShouldBindJSON(body); err != nil || body.Input == "" {
		c.JSON(400, ChatResponse{Error: "Invalid request"})
		return
	}

	res, err := s.openAi.SubmitChat(body.Input, GetAccessToken(c), c)
	if err != nil {
		if err == types.ErrAuthCode {
			c.JSON(401, ChatResponse{Error: "Unauthorized"})
			return

		} else {
			u, err := uuid.NewUUID()
			c.JSON(500, ChatResponse{Error: "Internal server error " + u.String()})
			logger.Logger().Err(err).Str("req_id", u.String()).Msg("Internal server error")
			return
		}
	}

	c.JSON(200, ChatResponse{Output: res})
}

func (s *serverImpl) RunAndBlock() error {
	bindAddr := conf.BindHost() + ":" + strconv.Itoa(conf.BindPort())
	logger.Logger().Info().Msgf("Listening on address %s", bindAddr)
	return s.engine.Run(bindAddr)
}
