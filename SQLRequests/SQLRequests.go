package SQLRequests

import (
	"WeatherBot/ApiRequests"
	ApiTokens "WeatherBot/key"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID         int              `json:"ID"`
	Username   string           `json:"Username"`
	ChatID     int64            `json:"ChatID"`
	City       ApiRequests.City `json:"City"`
	CallsLimit int              `json:"CallsLimit"`
	CallsLeft  int              `json:"CallsLeft"`
	HasSub     bool             `json:"HasSub"`
}
type Notification struct {
	ID           int
	chatid       int64
	weekDays     string
	time         string
	number       int
	wasSentToday bool
	isChanging   bool
}

func AddCity(city ApiRequests.City) {
	if CheckIfCityExistsInDB(city[0].LocalNames.Ru) {
		return
	}
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into cities(cityName, Lat, Lon, WasUpdated, TimeZone) values(?, ?, ?, ?, ?)",
		city[0].LocalNames.Ru, city[0].Lat, city[0].Lon, 0, ApiRequests.GetTimeZone(city))
	if err != nil {
		panic(err)
	}

}
func AddUser(chatid int64) {
	if CheckIfUserExists(chatid) {
		return
	}
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into users(chatid, relocationsLeft, relocationsTries, hasSub) values(?, 5, 5, 0)", chatid)
	if err != nil {
		panic(err)
	}
}
func CheckIfUserExists(chatid int64) bool {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRow("select exists(select * from Users where chatid = ?)", chatid)
	exist := false
	row.Scan(&exist)
	return exist
}
func CheckIfCityExistsInDB(cityName string) bool {
	fmt.Print(cityName)
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRow("select exists(select * from cities where cityName = ?)", cityName)
	exist := false
	row.Scan(&exist)
	return exist
}

func SetUserCity(cityName string, chatid int64) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update users set city = ? where chatid = ?", cityName, chatid)
	if err != nil {
		panic(err)
	}
}
func CheckWeatherActuality(cityName string) bool {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select wasUpdated from cities where cityName = ?", cityName)
	if err != nil {
		panic(err)
	}
	var a int
	for row.Next() {
		err := row.Scan(&a)
		if err != nil {
			panic(err)
		}
	}
	return a != 0
}
func UpdateWeather(cityName string) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var lat, lon float64
	row, err := db.Query("select lat, lon from cities where cityName = ?", cityName)
	if err != nil {
		panic(err)
	}
	for row.Next() {
		err := row.Scan(&lat, &lon)
		if err != nil {
			panic(err)
		}
	}
	weather, err := json.Marshal(ApiRequests.CheckTodaysWeather(lat, lon))
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("update cities set todaysWeather = ? where cityName = ?", weather, cityName)
	if err != nil {
		panic(err)
	}
	SetWeatherActualityTrue(cityName)
}
func SetWeatherActualityTrue(cityName string) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update cities set wasUpdated = 1 where cityName = ?", cityName)
	if err != nil {
		panic(err)
	}
}
func SetWeatherActualityFalse() {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	curtime := time.Now().UTC().Hour()
	_, err = db.Exec("update cities set wasUpdated = 0 where (timeZone + ? = 0 or timeZone + ? = 24)", curtime, curtime)
	if err != nil {
		panic(err)
	}
}

// вызывается в server при попытке запроса погоды, затем, если город указан, апдейтит погоду и возвращает юзеру результат
func GetUserCityName(chatid int64) string {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select city from users where chatid = ?", chatid)
	if err != nil {
		panic(err)
	}
	var cityName sql.NullString
	for row.Next() {
		err := row.Scan(&cityName)
		if err != nil {
			panic(err)
		}
	}
	return cityName.String
}
func ResetRelocationsTries() {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update users set relocationsTries = 5")
	if err != nil {
		panic(err)
	}
}
func ResetRelocationsLeft() {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update users set relocationsLeft = 5")
	if err != nil {
		panic(err)
	}
}
func LowerRelocationsLeft(chatid int64) {
	left, _ := GetRelocationsLeftnTries(chatid)
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update users set relocationsLeft = ?", left-1)
	if err != nil {
		panic(err)
	}
}
func LowerRelocationsTries(chatid int64) {
	_, tries := GetRelocationsLeftnTries(chatid)
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update users set relocationsLeft = ?", tries-1)
	if err != nil {
		panic(err)
	}
}
func GetRelocationsLeftnTries(chatid int64) (int, int) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select relocationsLeft, relocationsTries from users where chatid = ?", chatid)
	if err != nil {
		panic(err)
	}
	var left, tries int
	for row.Next() {
		err := row.Scan(&left, &tries)
		if err != nil {
			panic(err)
		}
	}
	return left, tries
}
func GetWeatherJSON(cityName string) ApiRequests.TodaysWeather {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRow("select todaysWeather from cities where cityName = ?", cityName)
	var weatherJSON []byte
	var weather ApiRequests.TodaysWeather
	row.Scan(&weatherJSON)
	err = json.Unmarshal(weatherJSON, &weather)
	if err != nil {
		panic(err)
	}
	return weather
}
func SetUserChangeStatus(chatid int64, status int) {
	fmt.Print(chatid)
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update users set changeStatus = ? where chatid = ?", status, chatid)
	if err != nil {
		panic(err)
	}
}
func GetUserChangeStatus(chatid int64) int {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select changeStatus from users where chatid = ?", chatid)
	if err != nil {
		panic(err)
	}
	var status int
	for row.Next() {
		err := row.Scan(&status)
		if err != nil {
			panic(err)
		}
	}
	return status
}
func GetUserUTC(chatid int64) int {
	cityName := GetUserCityName(chatid)
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select timeZone from cities where cityName = ?", cityName)
	if err != nil {
		panic(err)
	}
	var utc int
	for row.Next() {
		err := row.Scan(&utc)
		if err != nil {
			panic(err)
		}
	}
	return utc
}
func GetUserChangingNotification(chatid int64) Notification {
	notification := Notification{}
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select * from notifications where chatid = ? and isChanging = ?", chatid, true)
	if err != nil {
		panic(err)
	}
	var i int
	for row.Next() {
		err := row.Scan(&notification.ID, &notification.chatid, &notification.weekDays, &notification.time,
			&notification.number, &notification.wasSentToday, &notification.isChanging)
		if err != nil {
			panic(err)
		}
		i++
	}
	//проверить, если нет изменяемого оповещения
	return notification
}
func GetCurrentNotifications(chatid int64, day int, currentTime time.Time) []Notification { // неправильно - переписать(изменить как находится нужное время и переключить wasSentToday)
	formattedTime := currentTime.Format("03:04")
	notifications := []Notification{}
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select * from notifications where (chatid = ? and weekDays like '%%?%%' and time = ?)", chatid, strconv.Itoa(day), formattedTime)
	if err != nil {
		panic(err)
	}
	var i int
	for row.Next() {
		notifications = append(notifications, Notification{})
		err := row.Scan(&notifications[i].ID, &notifications[i].chatid, &notifications[i].weekDays, &notifications[i].time,
			&notifications[i].number, &notifications[i].wasSentToday, &notifications[i].isChanging)
		if err != nil {
			panic(err)
		}
		i++
	}
	return notifications
}
func GetNumberOfNotifications(chatid int64) int {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT count(*) FROM notifications where chatid = ?", chatid)
	if err != nil {
		panic(err)
	}
	var number int
	for row.Next() {
		err := row.Scan(&number)
		if err != nil {
			panic(err)
		}
	}
	return number
}
func DeleteNotification(chatid int64, number int) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("delete from notifications where chatid = ? and number = ?", chatid, number)
	if err != nil {
		panic(err)
	}
	ChangeNotificationNumbers(chatid, number)
}
func ChangeNotificationNumbers(chatid int64, number int) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("update notifications set number = (number - 1) where chatid = ? and number > ?", chatid, number)
	if err != nil {
		panic(err)
	}
}
func GetMaxNotificationNumber(chatid int64) int {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("SELECT MAX(number) FROM notifications where chatid = ?", chatid)
	if err != nil {
		panic(err)
	}
	var number int
	for row.Next() {
		err := row.Scan(&number)
		if err != nil {
			panic(err)
		}
	}
	return number
}
func AddNotification(chatid int64) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	number := GetMaxNotificationNumber(chatid) + 1
	_, err = db.Exec("insert into notifications(chatid, number) values(?, ?)", chatid, number)
	if err != nil {
		panic(err)
	}
}
func GetChangableNotificationWeekdays(chatid int64) string {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select weekDays from notifications where chatid = ? and isChanging = ?", chatid, true)
	if err != nil {
		panic(err)
	}
	var weekdays string
	for row.Next() {
		err := row.Scan(&weekdays)
		if err != nil {
			panic(err)
		}
	}
	return weekdays
}
func ChangeNotificationDay(chatid int64, weekdayint int) {
	weekdays := GetChangableNotificationWeekdays(chatid)
	weekday := strconv.Itoa(weekdayint)
	if !strings.Contains(weekdays, weekday) {
		weekdays += weekday
	} else {
		weekdays = weekdays[:strings.Index(weekdays, weekday)] + weekdays[strings.Index(weekdays, weekday)+1:]
	}
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("update notifications set weekdays = ? where chatid = ?, isChanging = ?", weekdays, chatid, true)
	if err != nil {
		panic(err)
	}
}
func UpdateIsChanging(chatid int64) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("update notifications set isChanging = ? where chatid = ? and isChanging = ?", false, chatid, true)
	if err != nil {
		panic(err)
	}
}
func SetIsChangingTrue(chatid int64, number int) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("update notifications set isChanging = ? where chatid = ? and number = ?", true, chatid, number)
	if err != nil {
		panic(err)
	}
}
func CheckIfGivenTimeIsOkay(givenTime string) bool {
	if _, err := time.Parse("03:04", givenTime); err != nil {
		return false
	}
	return true
}

//todo: обнулять переменную city/wasupdated в 00:01 по местному времени
//если пользователь запрашивает данные с city где wasupdate == false(0), то запрашивать нынешнюю погоду в данном городе
//может быть запретить менять город чаще n раз за месяц

/*
Список функций, которые надо реализовать(не факт что полный):
1)+ Проверить, какой город указан у полльзователя(sql запрос)
2)+ Есть ли город в БД(sql запрос)
3)+ Установить город пользователю(sql запрос)
4)+ Проверка реальности города(api request)
5)+ Добавление нового города в БД(sql запрос)
6)+ Проверка, что погода обновлена сегодня(sql запрос)
7)+ Узнать погоду на сегодня(api request)
8)+ Проверка, что время == XX:01
9)+ Обновление статуса на 0 для городов, у которых 12 ночи по местному времени(sql запрос)
9.1)+ Обновление relocationsLeft в начале месяца(не по местному)
9.2)+ Обновление relocationsTries в начале дня(не по местному)
10)+ Узнать местное время в новом(для БД) городе(api request)
11)+ Проверка на наличие пользователя в бд(при нажатии на /start)(sql запрос)
12)+ Добавление пользователя в БД(sql запрос)
13) Уменьшение relocationsLeft при успешной смене города
13.5) Уменьшение relocationsTries при не успешной смене города
14)+ Получение всех notification, нужные на данное время
15)+ Получение notifications пользователя (надо?)
15.5)+ Получение количества notifications пользователя
16)+ Удаление notification под литерой chatid notnumber
16.5)+ Изменение notnumber на -1 для всех notification где notnumber > notnumber удалённого уведомления
17)+ Добавление notification с нулевыми значениями
18)+ Изменение нужного notification
19) Проверка на наличие уведомления для пользователя на данное время и отправка уведомления
20) get user id message
*/
