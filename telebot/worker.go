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

		var showedJobAnnoucesIDs []uint

		for _, ja := range convertJobDataModelDBtoTG(a) {
			if err = ja.sentJobAnnounceToClient(ctx, ud.TgID, b); err != nil {
				logger.Error(err.Error())
				continue
			}

			showedJobAnnoucesIDs = append(showedJobAnnoucesIDs, ja.ItemID)
		}

		if len(showedJobAnnoucesIDs) != 0 {
			if err = bd.CreatePivotVacancyAnnouncesAndUserIds(showedJobAnnoucesIDs, uint(ud.TgID)); err != nil {
				logger.Error(err.Error())
				continue
			}
		}

	}
}
