package main

import (
	GisMeteoRequest "WeatherBot/GisMeteoRequests"
	ApiTokens "WeatherBot/key"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commandKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Погода сейчас", "Узнаём погоду..."),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Погода сегодня", "Узнаём погоду на сегодня..."),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить город", "Меняем город..."),
	),
)
var problem string
var currentWeather GisMeteoRequest.CurrentWeather
var todaysWeather GisMeteoRequest.TodaysWeather
var waitingForCity bool
var city GisMeteoRequest.City

func main() {
	bot, err := tgbotapi.NewBotAPI(ApiTokens.BotToken)
	if err != nil {
		fmt.Println(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я - WeatherBot! Я помогу тебе узнать погоду в твоём городе, главное не забудь указать его!")
					msg.ReplyMarkup = commandKeyboard
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
				}
			} else if update.Message.Text != "" {
				if waitingForCity {
					city, problem = GisMeteoRequest.UpdateCity(update.Message.Text)
					if problem != "" || len(city) == 0 {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Город с таким названием не найден, перепроверьте его или напишите другое.")
						city = nil
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Город успешно изменён")
						msg.ReplyMarkup = commandKeyboard
						bot.Send(msg)
						waitingForCity = false
					}

				}
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
			switch update.CallbackQuery.Data {
			case "Узнаём погоду...":
				if city == nil {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Город не указан, укажите город")
					bot.Send(msg)
					waitingForCity = true
					continue
				}
				currentWeather, problem = GisMeteoRequest.CheckCurrentWeather(city)
				if problem != "" {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, problem)
					bot.Send(msg)
					continue
				}
				answer := fmt.Sprintf("Погода в %v:\nСейчас %v°C\nСкорость ветра %vм/c\nВлажность %v%%\nОсадки %vсм", city[0].LocalNames.Ru,
					currentWeather.Data.Temperature.Air.C, currentWeather.Data.Wind.Speed.MS,
					currentWeather.Data.Humidity.Percent, currentWeather.Data.Precipitation.Amount)
				answer += "\n\n\nПодробнее <a href=\"https://www.gismeteo.ru\">здесь</a>"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = commandKeyboard
				bot.Send(msg)
			case "Узнаём погоду на сегодня...":
				if city == nil {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Город не указан, укажите город")
					bot.Send(msg)
					waitingForCity = true
					continue
				}
				todaysWeather, problem = GisMeteoRequest.CheckTodaysWeather(city)
				answer := fmt.Sprintf("Погода в %v на сегодня:\n", city[0].LocalNames.Ru)
				for i, period := range todaysWeather.Data {
					answer += fmt.Sprintf("В %v:00 ожидается %v°C\nСкорость ветра %v метров в секунду\nВлажность %v%%\nКоличество осадков около %vмм \n\n", i*3,
						period.Temperature.Air.C, period.Wind.Speed.MS, period.Humidity.Percent, period.Precipitation.Amount)
					if i >= 8 {
						break
					}
				}
				answer += "\nПодробнее <a href=\"https://www.gismeteo.ru\">здесь</a>"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
				msg.ReplyMarkup = commandKeyboard
				msg.ParseMode = "HTML"
				bot.Send(msg)
			case "Меняем город...":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Напишите название города")
				bot.Send(msg)
				waitingForCity = true
			default:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Что-то пошло не так")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}

		}

	}
}

//TODO1: Добавить базу данных с пользователями(айди, никнейм, город, ежедневное оповещение о погоде),
//сегодняшней и нынешней погодой в городах(название города, байты(предположительно, не изменяются в течение дня))
//TODO2: Добавить возможность присылать координаты заместо города
//TODO1.5:Сделать проверку погоду на сегодня раз в 24 часа в 00:00 по местному времени
