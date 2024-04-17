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
var frequency = 4

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
				answer := fmt.Sprintf("Сейчас %v°C, ощущается как %v°C\nСкорость ветра %v метров в секунду\nВлажность %v%%",
					currentWeather.Current.TempC, currentWeather.Current.FeelslikeC, fmt.Sprintf("%.1f", currentWeather.Current.WindKph/3.6), currentWeather.Current.Humidity)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
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
				answer := ""
				for i := 0; i < 24; i += frequency {
					answer += fmt.Sprintf("В %v:00 ожидается %v°C, ощущается как %v°C\nСкорость ветра %v метров в секунду\nВлажность %v%%\n\n", i, //TODO: дождик/осадки
						todaysWeather.Forecast.Forecastday[0].Hour[i].TempC, todaysWeather.Forecast.Forecastday[0].Hour[i].FeelslikeC,
						fmt.Sprintf("%.1f", todaysWeather.Forecast.Forecastday[0].Hour[i].WindKph/3.6), todaysWeather.Forecast.Forecastday[0].Hour[i].Humidity)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
				msg.ReplyMarkup = commandKeyboard
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
