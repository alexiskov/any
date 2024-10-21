package telebot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"vacancydealer/bd"
	"vacancydealer/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// -------------------------------------------------------------------------------------->>>HANDLERS------------------------------------------------------------------------------
// Client message processing handler
func textHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgUID := update.Message.From.ID

	switch update.Message.Text {
	default:
		if u, ok := UserStates[tgUID]; ok {
			defer delete(UserStates, tgUID)

			switch u.State {
			case 1:
				u.User.Vacancy = update.Message.Text
				if err := u.User.convertUserModelTGtoDB().Update(); err != nil {
					logger.Error(err.Error())
					return
				}

				if err := sentUserDataToClient(ctx, tgUID, b); err != nil {
					logger.Error(err.Error())
					return
				}
			case 21:
				cities, err := bd.FindCitiesByName(update.Message.Text)
				if err != nil {
					logger.Error(err.Error())
				}

				buttonsData := make([][2]string, 0)
				for _, city := range cities {
					buttonsData = append(buttonsData, [2]string{city.Name, "?setLocation:" + strconv.Itoa(int(city.ID))})
				}

				msgParams := &bot.SendMessageParams{
					ChatID:    tgUID,
					ParseMode: models.ParseModeHTML,
				}
				if len(buttonsData) != 0 && len(buttonsData) < 30 {
					msgParams.Text = "<b>Уточним локацию</b>\n\nНажми нужную кнопку."
					msgParams.ReplyMarkup = &models.InlineKeyboardMarkup{InlineKeyboard: linesButtonGenerate(buttonsData)}
				} else {
					msgParams.Text = "<b>Уточним локацию</b>\n\nНет результатов, пожалуйста уточните название населенного пункта."
					msgParams.ReplyMarkup = &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: "не имеет значения", CallbackData: "?setLocation:0"}}}}
				}
				_, err = b.SendMessage(ctx, msgParams)
				if err != nil {
					logger.Error(err.Error())
				}
			case 22:
				regions, err := bd.FindRegionByName(update.Message.Text)
				if err != nil {
					logger.Error(err.Error())
				}

				buttonsData := make([][2]string, 0)
				for _, city := range regions {
					buttonsData = append(buttonsData, [2]string{city.Name, "?setLocation:" + strconv.Itoa(int(city.ID))})
				}

				msgParams := &bot.SendMessageParams{
					ChatID:    tgUID,
					ParseMode: models.ParseModeHTML,
				}
				if len(buttonsData) != 0 && len(buttonsData) < 30 {
					msgParams.Text = "<b>Уточним локацию</b>\n\nНажми нужную кнопку."
					msgParams.ReplyMarkup = &models.InlineKeyboardMarkup{InlineKeyboard: linesButtonGenerate(buttonsData)}
				} else {
					msgParams.Text = "<b>Уточним локацию</b>\n\nНет результатов, пожалуйста уточните название населенного пункта."
					msgParams.ReplyMarkup = &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: "не имеет значения", CallbackData: "?setLocation:0"}}}}
				}
				_, err = b.SendMessage(ctx, msgParams)
				if err != nil {
					logger.Error(err.Error())
				}
			case 3:

				exp, err := strconv.Atoi(update.Message.Text)
				if err != nil {
					logger.Error(fmt.Errorf("input experience value parsing error: %w", err).Error())
					return
				}
				u.User.ExperienceYears = exp
				if err := u.User.convertUserModelTGtoDB().Update(); err != nil {
					logger.Error(err.Error())
					return
				}

				if err := sentUserDataToClient(ctx, tgUID, b); err != nil {
					logger.Error(err.Error())
					return
				}
			}
		} else {
			if err := sentUserDataToClient(ctx, tgUID, b); err != nil {
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
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      tgUID,
			ParseMode:   models.ParseModeHTML,
			Text:        "<b>Что изменим?</b>\n\nНажми нужную кнопку.",
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: linesButtonGenerate([][2]string{{"профессия", "#changeVacancyName"}, {"регион", "#changeLocation"}, {"опыт работы", "#changeExperience"}, {"график работы", "#changeSchedule"}})},
		})
		if err != nil {
			logger.Error(fmt.Errorf("filter write command handler error^ %w", err).Error())
		}
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
			ChatID:      tgUID,
			ParseMode:   models.ParseModeHTML,
			Text:        "<b>Замена региона поиска вакансии</b>\n\nУточнить локацию поиска до:",
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: linesButtonGenerate([][2]string{{"страны", "#changeCountry"}, {"региона", "#changeRegion"}, {"населенного пункта", "#changeCity"}, {"не имеет значения", "?setLocation:0"}})},
		})
		if err != nil {
			logger.Error(fmt.Errorf("change city name function, to user %d have a error: %w", tgUID, err).Error())
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
	case "#changeRegion":
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    tgUID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>Укажите регион/область</b>\n\nВведите название региона/области:",
		})
		if err != nil {
			logger.Error(fmt.Errorf("change region name function, to user %d have a error: %w", tgUID, err).Error())
			return
		}

		state := UserStateData{State: 22, Date: time.Now()}
		state.User, err = findRegisterUser(tgUID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		UserStates[tgUID] = state
	case "#changeCountry":
		countries, err := bd.FindCountries()
		if err != nil {
			logger.Error(err.Error())
		}

		buttonsData := make([][2]string, 0)
		for _, city := range countries {
			buttonsData = append(buttonsData, [2]string{city.Name, "?setLocation:" + strconv.Itoa(int(city.ID))})
		}

		msgParams := &bot.SendMessageParams{
			ChatID:    tgUID,
			ParseMode: models.ParseModeHTML,
		}
		if len(buttonsData) != 0 && len(buttonsData) < 30 {
			msgParams.Text = "<b>Уточним локацию</b>\n\nНажми нужную кнопку."
			msgParams.ReplyMarkup = &models.InlineKeyboardMarkup{InlineKeyboard: linesButtonGenerate(buttonsData)}
		} else {
			msgParams.Text = "<b>Уточним локацию</b>\n\nНет результатов, пожалуйста уточните название населенного пункта."
			msgParams.ReplyMarkup = &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{{{Text: "не имеет значения", CallbackData: "?setLocation:0"}}}}
		}

		_, err = b.SendMessage(ctx, msgParams)
		if err != nil {
			logger.Error(err.Error())
		}
	case "#changeExperience":
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    tgUID,
			ParseMode: models.ParseModeHTML,
			Text:      "<b>Опыт работы</b>\n\nУкажите в годах, Ваш опыт в искомой сфере - числом\n <u>пример:</u> 12"},
		)
		if err != nil {
			logger.Error(fmt.Errorf("change vacancy name function, to user %d have a error: %w", tgUID, err).Error())
			return
		}

		state := UserStateData{State: 3, Date: time.Now()}
		state.User, err = findRegisterUser(tgUID)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		UserStates[tgUID] = state
	case "#changeSchedule":

		sch, err := bd.GetSchedules("")
		if err != nil {
			logger.Error(err.Error())
			return
		}
		schedulesButtonsData := make([][2]string, 0)

		for _, s := range sch {
			schedulesButtonsData = append(schedulesButtonsData, [2]string{s.Name, "?changeSched:" + s.HhID})
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      tgUID,
			ParseMode:   models.ParseModeHTML,
			Text:        "<b>График работы</b>\n\nВыберите график работы по искомой вакансии",
			ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: linesButtonGenerate(schedulesButtonsData)}},
		)
		if err != nil {
			logger.Error(fmt.Errorf("change vacancy name function, to user %d have a error: %w", tgUID, err).Error())
			return
		}
	}

}

// SetLocation Handler
func locationSetter(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgUID := update.CallbackQuery.From.ID
	locationID, err := strconv.Atoi(strings.Trim(update.CallbackQuery.Data, "?setLocation:"))
	if err != nil {
		err = fmt.Errorf("incomming callbackData of region id parsing error: %w", err)
	}
	u, err := findRegisterUser(tgUID)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	sqluser := u.convertUserModelTGtoDB()
	sqluser.Location = uint(locationID)
	if err = sqluser.UpdateLocation(); err != nil {
		logger.Error(err.Error())
	}

	sentUserDataToClient(ctx, tgUID, b)
}

// Set schedule handler
func scheduleSetter(ctx context.Context, b *bot.Bot, update *models.Update) {
	tgUID := update.CallbackQuery.From.ID
	u, err := findRegisterUser(tgUID)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	sqluser := u.convertUserModelTGtoDB()
	sqluser.Schedule = strings.Replace(update.CallbackQuery.Data, "?changeSched:", "", 1)
	if err = sqluser.UpdateSchedule(); err != nil {
		logger.Error(err.Error())
	}
	if err = sentUserDataToClient(ctx, tgUID, b); err != nil {
		logger.Error(err.Error())
		return
	}
}

// --------------------------------------------------------------------------------------<<<HANDLERS------------------------------------------------------------------------------

// ------------------------------------------------------------------------->>>BUTTON GENERATOR---------------------------------------------------------------
func linesButtonGenerate(buttonsData [][2]string) (inlineButtons [][]models.InlineKeyboardButton) {
	for _, lb := range buttonsData {
		inlineButtons = append(inlineButtons, []models.InlineKeyboardButton{{Text: lb[0], CallbackData: lb[1]}})
	}
	return
}

// -------------------------------------------------------------------------<<<BUTTON GENERATOR---------------------------------------------------------------
