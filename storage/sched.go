package storage

import (
	"github.com/collinzh/chatbot-messenger/logger"
	"github.com/collinzh/chatbot-messenger/util"
	"time"
)

func ScheduleMessagePruning(sto Storage, delay time.Duration) {
	util.ScheduleTask(func() {
		logger.Logger().Debug().Msg("Pruning chats")
		if err := sto.Prune(); err != nil {
			logger.Logger().Err(err).Msg("Error pruning storage")
		}
	}, delay)
}
