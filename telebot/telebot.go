package telebot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"vacancydealer/bd"
	"vacancydealer/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type (
	UserData struct {
		TgID            int64
		Vacancy         string
		City            string
		Schedule        string
		ExperienceYears int
	}
)

var (
	UserStates map[string]string
)

// Start tgelegram-Bot worker
func Run(tgAPI string) (err error) {
	UserStates = make(map[string]string, 100)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMessageTextHandler("", bot.MatchTypeContains, textHandler),
		bot.WithCallbackQueryDataHandler("#", bot.MatchTypePrefix, callbackProcessing),
	}

	b, err := bot.New(tgAPI, opts...)
	if err != nil {
		return
	}

	b.Start(ctx)

	return nil
}

// Client message processing handler
func textHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	u, err := findRegisterUser(update.Message.From.ID)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	switch update.Message.Text {
	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.From.ID,
			ParseMode: models.ParseModeHTML,
			Text:      fmt.Sprintf("<b> <u>Поиск вакансий</u> </b>\n\n<b>Профессия: </b><i> %s</i>\n<b>Город: </b><i> %s</i>\n<b>Опыт работы(лет): </b> %d\n<b>График работы: </b> <i> %s</i>", u.Vacancy, u.City, u.ExperienceYears, u.Schedule),
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "редактировать", CallbackData: "#vacFilterWritePlease"}},
			}},
		})
	}
}

// Client callback handler
func callbackProcessing(ctx context.Context, b *bot.Bot, update *models.Update) {
	u, err := findRegisterUser(update.CallbackQuery.From.ID)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	switch update.CallbackQuery.Data {
	case "#vacFilterWritePlease":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>Что изменим?</b>\n\nНажми нужную кнопку.",
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "профессия", CallbackData: "#changeVacancyName"}},
				{{Text: "город", CallbackData: "#changeCity"}},
				{{Text: "опыт работы", CallbackData: "#changeExperience"}},
				{{Text: "график работы", CallbackData: "#changeSchedule"}},
			}},
		})
	case "changeVacancyName":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.CallbackQuery.From.ID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>назвние вакансие?</b>\n\nНазвание вакансии не обязательно должно быть полным. Поиск происходит по совпадению ключевых слов в названии вакансии. Допустимо указать одно слово в вакансии или полное название. Важно понимать, что работодатель указывает произвольное название.\n\nвведи ключевое слово для поиска вакансии"})
	}
}

// Finding
func findRegisterUser(tgID int64) (ud UserData, err error) {
	sqludata, err := bd.FindOrCreateUser(tgID)
	if err != nil {
		return
	}
	ud = UserData{TgID: sqludata.TgID, Vacancy: sqludata.VacancyName, ExperienceYears: sqludata.ExperienceYear}

	if sqludata.City == "" {
		ud.City = "не имеет значения"
	} else {
		ud.City = sqludata.City
	}

	switch sqludata.Schedule {
	case 1:
		ud.Schedule = "удаленная работа"
	default:
		ud.Schedule = "не выбрано"
	}

	if sqludata.VacancyName == "" {
		ud.Vacancy = "не указано"
	}
	return
}

func (ud UserData) UpdateSearchFilter() {

}
