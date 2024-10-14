package telebot

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type (
	UserData struct {
		TgID            uint64
		Vacancy         string
		City            string
		Schedule        string
		ExperienceYears int8
	}
)

func Run(tgAPI string) (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(tgAPI, opts...)
	if err != nil {
		return
	}

	b.Start(ctx)

	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	ud := findRegisterUser(update.Message.From.ID)

}

func findRegisterUser(tgID int64) (ud UserData) {
	return
}
