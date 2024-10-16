package telebot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
	"vacancydealer/bd"
	"vacancydealer/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type (
	UserData struct {
		TgID            int64
		Vacancy         string
		Location        string
		Schedule        string
		ExperienceYears int
	}

	UserStateData struct {
		State uint8
		User  UserData
		Date  time.Time
	}
)

var (
	UserStates map[int64]UserStateData
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
	}

	b, err := bot.New(tgAPI, opts...)
	if err != nil {
		return
	}

	b.Start(ctx)

	return nil
}

// -------------------------------------------------------------------------------------->>>HANDLERS------------------------------------------------------------------------------
// Client message processing handler
func textHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgUID := update.Message.From.ID

	switch update.Message.Text {
	default:
		if u, ok := UserStates[tgUID]; ok {
			switch u.State {
			case 1:
				defer delete(UserStates, tgUID)

				u.User.Vacancy = update.Message.Text
				if err := u.User.convertUserModelTGtoDB().Update(); err != nil {
					logger.Error(err.Error())
					return
				}

				if err := showUserData(ctx, tgUID, b); err != nil {
					logger.Error(err.Error())
					return
				}
			case 21:
				defer delete(UserStates, tgUID)

				cities, err := bd.FindCitiesByName(update.Message.Text)
				if err != nil {
					logger.Error(err.Error())
				}

				buttonsData := make([][2]string, 0)
				for _, city := range cities {
					buttonsData = append(buttonsData, [2]string{city.Name, "?setLocation:" + strconv.Itoa(int(city.ID))})
				}

				_, err = b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      tgUID,
					ParseMode:   models.ParseModeHTML,
					Text:        "<b>Уточним локацию</b>\n\nНажми нужную кнопку.",
					ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: linesButtonGenerate(buttonsData)},
				})
				if err != nil {
					logger.Error(err.Error())
				}

			}
		} else {
			if err := showUserData(ctx, tgUID, b); err != nil {
				logger.Error(err.Error())
				return
			}
		}

	}
}

// Client callback handler
func callbackProcessing(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgUID := update.CallbackQuery.From.ID

	switch update.CallbackQuery.Data {
	case "#vacFilterWritePlease":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    tgUID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>Что изменим?</b>\n\nНажми нужную кнопку.",
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "профессия", CallbackData: "#changeVacancyName"}},
				{{Text: "регион", CallbackData: "#changeLocation"}},
				{{Text: "опыт работы", CallbackData: "#changeExperience"}},
				{{Text: "график работы", CallbackData: "#changeSchedule"}},
			}},
		})
	case "#changeVacancyName":
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    tgUID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>назвние вакансии</b>\n\nНазвание вакансии не обязательно должно быть полным. Поиск происходит по совпадению ключевых слов в названии вакансии. Допустимо указать одно слово в вакансии или более. Важно понимать, что работодатель указывает произвольное название.\n\nвведи ключевое слово для поиска по названию вакансии"},
		)
		if err != nil {
			logger.Error(fmt.Errorf("change vacancy name function, to user %d have a error: %w", tgUID, err).Error())
			return
		}

		state := UserStateData{State: 1, Date: time.Now()}
		state.User, err = findRegisterUser(tgUID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		UserStates[tgUID] = state
	case "#changeLocation":
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    tgUID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>Замена региона поиска вакансии</b>\n\nУточнить локацию поиска до:",
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "страны", CallbackData: "#changeCountry"}},
				{{Text: "региона", CallbackData: "#changeRegion"}},
				{{Text: "населенного пункта", CallbackData: "#changeCity"}},
			}},
		})
		if err != nil {
			logger.Error(fmt.Errorf("change vacancy name function, to user %d have a error: %w", tgUID, err).Error())
			return
		}
	case "#changeCity":
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    tgUID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>Укажите населенный пункт</b>\n\nВведите название населенного пункта:",
		})
		if err != nil {
			logger.Error(fmt.Errorf("change vacancy name function, to user %d have a error: %w", tgUID, err).Error())
			return
		}

		state := UserStateData{State: 21, Date: time.Now()}
		state.User, err = findRegisterUser(tgUID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		UserStates[tgUID] = state
	}

}

// SetLocation Handler
func locationSetter(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("sadasdasdasdsadsdsdsdsdsdsdsdsdsd")
	tgUID := update.CallbackQuery.From.ID
	locationID, err := strconv.Atoi(strings.Trim(update.CallbackQuery.Data, "?setLocation:"))
	if err != nil {
		err = fmt.Errorf("incomming callbackData of region id parsing error: %w", err)
	}
	u, err := findRegisterUser(tgUID)
	if err != nil {
		return
	}

	sqluser := u.convertUserModelTGtoDB()
	sqluser.Location = uint(locationID)
	sqluser.Update()
}

// --------------------------------------------------------------------------------------<<<HANDLERS------------------------------------------------------------------------------

// Find or Write user on db
func findRegisterUser(tgID int64) (ud UserData, err error) {
	sqludata, err := bd.FindOrCreateUser(tgID)
	if err != nil {
		return
	}
	return convertUserModelDBtoTG(sqludata), nil
}

// UserData response to tg-chat sent
func showUserData(ctx context.Context, tgID int64, b *bot.Bot) (err error) {
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

//-------------------------------------MODEL CONVERTERS-----------------------------------

func (ud UserData) convertUserModelTGtoDB() (sqluser bd.UserData) {
	if ud.Location == "не имеет значения" {
		sqluser.Location = 0
	} else {
		//sqluser.Location = ud.Location
	}

	switch ud.Schedule {
	case "удаленная работа":
		sqluser.Schedule = 1
	default:
		sqluser.Schedule = 0
	}

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
		//ud.Location = sqluser.Location
	}

	switch sqluser.Schedule {
	case 1:
		ud.Schedule = "удаленная работа"
	default:
		ud.Schedule = "не выбрано"
	}

	if sqluser.VacancyName == "" {
		ud.Vacancy = "не указано"
	}
	return
}

// ------------------------------------------------------------------------->>>BUTTON GENERATOR---------------------------------------------------------------
func linesButtonGenerate(buttonsData [][2]string) (inlineButtons [][]models.InlineKeyboardButton) {
	for _, lb := range buttonsData {
		inlineButtons = append(inlineButtons, []models.InlineKeyboardButton{{Text: lb[0], CallbackData: lb[1]}})
	}
	return
}
