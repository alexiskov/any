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

var (
	UserStates     map[int64]UserStateData
	SCHEDULE_TYPES = []ScheduleType{{"удаленная работа", 1}, {"полная занятость", 2}}
)

// Start tgelegram-Bot worker
func Run(tgAPI string) (err error) {
	UserStates = make(map[int64]UserStateData, 100)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMessageTextHandler("", bot.MatchTypeContains, textHandler),
		bot.WithCallbackQueryDataHandler("#", bot.MatchTypePrefix, callbackProcessing),
		bot.WithCallbackQueryDataHandler("?setLocation:", bot.MatchTypePrefix, locationSetter),
		bot.WithCallbackQueryDataHandler("?changeSched:", bot.MatchTypePrefix, scheduleSetter),
	}

	b, err := bot.New(tgAPI, opts...)
	if err != nil {
		return
	}

	b.Start(ctx)

	return nil
}

// Find or Write user on db
func findRegisterUser(tgID int64) (ud UserData, err error) {
	sqludata, err := bd.FindOrCreateUser(tgID)
	if err != nil {
		return
	}
	return convertUserModelDBtoTG(sqludata), nil
}

// UserData response to tg-chat sent
func sentUserDataToClient(ctx context.Context, tgID int64, b *bot.Bot) (err error) {
	ud, err := findRegisterUser(tgID)
	if err != nil {
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    ud.TgID,
		ParseMode: models.ParseModeHTML,
		Text:      fmt.Sprintf("<b> <u>Поиск вакансий</u> </b>\n\n<b>Профессия: </b><i> %s</i>\n<b>Регион: </b><i> %s</i>\n<b>Опыт работы(лет): </b> %d\n<b>График работы: </b> <i> %s</i>", ud.Vacancy, ud.Location, ud.ExperienceYears, ud.Schedule),
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "редактировать", CallbackData: "#vacFilterWritePlease"}},
		}},
	})
	if err != nil {
		err = fmt.Errorf("UserData show error: %w", err)
		return
	}
	return nil
}

func (ja JobAnnounce) sentJobAnnounceToClient(ctx context.Context, tgID int64, b *bot.Bot) (err error) {
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    tgID,
		ParseMode: models.ParseModeHTML,
		Text:      fmt.Sprintf("<b> <u>%s</u> </b>\n\n<b>Требуемый опыт: </b><i> %s</i>\n<b>Зп указана до уплаты налогов: </b>%t\n<b>Размер ЗП:</b>%.2f - %.2f%s\n<b>Графика работы: </b>%s", ja.Name, ja.Expierence, ja.SalaryGross, ja.SalaryFrom, ja.SalaryTo, ja.SalaryCurrency, ja.Schedule),
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "источник", URL: ja.Link}},
		}},
	})
	if err != nil {
		err = fmt.Errorf("sentJobAnnounceTo client error: %w", err)
		return
	}
	return nil
}

//------------------------------------->>>MODEL CONVERTERS-----------------------------------

func (ud UserData) convertUserModelTGtoDB() (sqluser bd.UserData) {
	sqluser.TgID = ud.TgID
	sqluser.VacancyName = ud.Vacancy
	sqluser.ExperienceYear = ud.ExperienceYears
	return
}

func convertUserModelDBtoTG(sqluser bd.UserData) (ud UserData) {
	ud = UserData{TgID: sqluser.TgID, Vacancy: sqluser.VacancyName, ExperienceYears: sqluser.ExperienceYear}

	if sqluser.Location == 0 {
		ud.Location = "не имеет значения"
	} else {
		loc, err := bd.FindLocByID(sqluser.Location)
		if err != nil {
			logger.Error(err.Error())
			ud.Location = "не имеет значения"
		} else {
			ud.Location = loc
		}
	}

	res, _ := bd.GetSchedules(sqluser.Schedule)
	ud.Schedule = res[0].Name

	if sqluser.VacancyName == "" {
		ud.Vacancy = "не указано"
	}
	return
}

func convertJobDataModelDBtoTG(jobSQLdata []bd.JobAnnounce) (ja []JobAnnounce) {
	for _, sj := range jobSQLdata {
		ja = append(ja, JobAnnounce{Name: sj.Name, Expierence: sj.Expierence, SalaryGross: sj.SalaryGross, SalaryFrom: sj.SalaryFrom, SalaryTo: sj.SalaryTo, SalaryCurrency: sj.SalaryCurrency, PublishedAt: sj.PublishedAt, Schedule: sj.Schedule, Requirement: sj.Requirement, Responsebility: sj.Responsebility, Link: sj.Link})
	}
	return
}

//-------------------------------------<<<MODEL CONVERTERS-----------------------------------
