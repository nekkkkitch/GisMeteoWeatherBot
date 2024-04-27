package main

import (
	"WeatherBot/ApiRequests"
	"WeatherBot/SQLRequests"
	ApiTokens "WeatherBot/key"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var menuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
var cancelCityKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отменить смену города", "Отмена"),
	))

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
					SQLRequests.AddUser(update.Message.Chat.ID)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я - WeatherBot! Я помогу тебе узнать погоду в твоём городе, главное не забудь указать его!")
					msg.ReplyMarkup = menuKeyboard
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
				}
			} else if update.Message.Text != "" {
				if SQLRequests.GetUserChangeStatus(update.Message.Chat.ID) == 1 {
					if problem := TryToSetUserCity(update.Message.Text, (update.Message.Chat.ID)); problem != "" {
						message := ""
						switch problem {
						case "TOWNNOTEXIST":
							message = "Город с таким названием не найден, перепроверьте его или напишите другое."
						case "TOWNINDICATED":
							message = "Этот город уже указан."
						}
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
						msg.ReplyMarkup = cancelCityKeyboard
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Город успешно изменён")
						msg.ReplyMarkup = menuKeyboard
						bot.Send(msg)
						SQLRequests.SetUserChangeStatus(update.Message.Chat.ID, 0)
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
				todaysWeather, problem := TryToGetWeather(update.CallbackQuery.Message.Chat.ID)
				if problem == "MISSINGTOWN" {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Город не указан")
					bot.Send(msg)
					break
				}
				utc := SQLRequests.GetUserUTC(update.CallbackQuery.Message.Chat.ID)
				currentHour := time.Now().UTC().Hour()
				currentHour += utc
				if currentHour%3 == 1 {
					currentHour -= 1
				} else if currentHour%3 == 2 {
					currentHour += 1
				}
				currentWeather := todaysWeather.Data[currentHour]
				answer := fmt.Sprintf("Погода в %v:\nСейчас %v°C\nСкорость ветра %vм/c\nВлажность %v%%\nОсадки %vмм",
					SQLRequests.GetUserCityName(update.CallbackQuery.Message.Chat.ID),
					currentWeather.Temperature.Air.C, currentWeather.Wind.Speed.MS,
					currentWeather.Humidity.Percent, currentWeather.Precipitation.Amount)
				answer += "\n\n\nПодробнее <a href=\"https://www.gismeteo.ru\">здесь</a>"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = menuKeyboard
				bot.Send(msg)
			case "Узнаём погоду на сегодня...":
				todaysWeather, problem := TryToGetWeather(update.CallbackQuery.Message.Chat.ID)
				if problem == "MISSINGTOWN" {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Город не указан")
					bot.Send(msg)
					break
				}

				answer := fmt.Sprintf("Погода в %v на сегодня:\n", SQLRequests.GetUserCityName(update.CallbackQuery.Message.Chat.ID))
				for i, period := range todaysWeather.Data {
					answer += fmt.Sprintf("В %v:00 ожидается %v°C\nСкорость ветра %v метров в секунду\nВлажность %v%%\nКоличество осадков около %vмм \n\n", i*3,
						period.Temperature.Air.C, period.Wind.Speed.MS, period.Humidity.Percent, period.Precipitation.Amount)
					if i >= 8 {
						break
					}
				}
				answer += "\nПодробнее <a href=\"https://www.gismeteo.ru\">здесь</a>"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, answer)
				msg.ReplyMarkup = menuKeyboard
				msg.ParseMode = "HTML"
				bot.Send(msg)
			case "Меняем город...":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Напишите название города")
				msg.ReplyMarkup = cancelCityKeyboard
				bot.Send(msg)
				SQLRequests.SetUserChangeStatus(update.CallbackQuery.Message.Chat.ID, 1)
			case "Отмена":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Что хотите сделать теперь?")
				msg.ReplyMarkup = menuKeyboard
				bot.Send(msg)
				SQLRequests.SetUserChangeStatus(update.CallbackQuery.Message.Chat.ID, 0)
			default:
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Что-то пошло не так")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}

		}

	}
}
func CheckTime() {
	if time.Now().Minute() == 1 {
		SQLRequests.SetWeatherActualityFalse()
	}
	if time.Now().Hour() == 0 {
		if time.Now().Day() == 1 {
			SQLRequests.ResetRelocationsLeft()
		}
		SQLRequests.ResetRelocationsTries()
	}
}
func TryToSetUserCity(cityName string, chatid int64) string {
	if cityName == SQLRequests.GetUserCityName(chatid) {
		return "TOWNINDICATED"
	}
	if SQLRequests.CheckIfCityExistsInDB(cityName) {
		SQLRequests.SetUserCity(cityName, chatid)
		return "Учпечно"
	}
	city, err := ApiRequests.CheckIfCityIsReal(cityName)
	if !err {
		SQLRequests.AddCity(city)
		SQLRequests.SetUserCity(cityName, chatid)
		SQLRequests.SetUserChangeStatus(chatid, 0)
		return "Учпечно"
	}
	return "TOWNNOTEXIST"
}
func TryToGetWeather(chatid int64) (ApiRequests.TodaysWeather, string) {
	var weather ApiRequests.TodaysWeather
	cityName := SQLRequests.GetUserCityName(chatid)
	if cityName == "" {
		return weather, "MISSINGTOWN"
	}
	if !SQLRequests.CheckWeatherActuality(cityName) {
		SQLRequests.UpdateWeather(cityName)
	}
	weather = SQLRequests.GetWeatherJSON(cityName)
	return weather, ""
}

//TODO: Добавить возможность присылать координаты заместо города
//TODO2: сделать красивенько чтобы сообщения присылались(обновление последнего сообщения)
