package telebot

import (
	"context"
	"vacancydealer/bd"
	"vacancydealer/logger"

	"github.com/go-telegram/bot"
)

func StartWorker(ctx context.Context, b *bot.Bot) {
	uds, err := bd.GetAllUserData()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, ud := range uds {
		a, err := ud.GetJobAnnounces()
		if err != nil {
			logger.Error(err.Error())
			continue
		}

		for _, ja := range convertJobDataModelDBtoTG(a) {
			if err = ja.sentJobAnnounceToClient(ctx, ud.TgID, b); err != nil {
				logger.Error(err.Error())
				continue
			}
		}
	}
}
