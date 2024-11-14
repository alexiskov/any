package telebot

import (
	"context"
	"time"
	"vacancydealer/bd"
	"vacancydealer/logger"

	"github.com/go-telegram/bot"
)

// Automatic worker
// New vacancieAnnounces to user sent
func StartWorker(ctx context.Context, b *bot.Bot) {
	for {
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

		if len(uds) != 0 {
			time.Sleep(time.Duration(10800/len(uds)) * time.Second) //period
		}

	}

}
