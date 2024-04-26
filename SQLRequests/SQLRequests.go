package SQLRequests

import (
	"WeatherBot/ApiRequests"
	ApiTokens "WeatherBot/key"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID               int              `json:"ID"`
	Username         string           `json:"Username"`
	ChatID           int              `json:"ChatID"`
	City             ApiRequests.City `json:"City"`
	TimeNotification time.Time        `json:"TimeNotification"`
	CallsLimit       int              `json:"CallsLimit"`
	CallsLeft        int              `json:"CallsLeft"`
	HasSub           bool             `json:"HasSub"`
}
type City struct {
	CityName       string                     `json:"CityName"`
	WeatherToday   ApiRequests.TodaysWeather  `json:"WeatherToday"`
	WeatherCurrent ApiRequests.CurrentWeather `json:"WeatherCurrent"`
	Lat            float64                    `json:"Lat"`
	Lon            float64                    `json:"Lon"`
	WasUpdated     bool                       `json:"WasUpdated"`
	TimeZone       int                        `json:"TimeZone"`
}
type TimeZone struct {
	Status           string      `json:"status"`
	Message          string      `json:"message"`
	CountryCode      string      `json:"countryCode"`
	CountryName      string      `json:"countryName"`
	RegionName       interface{} `json:"regionName"`
	CityName         string      `json:"cityName"`
	ZoneName         string      `json:"zoneName"`
	Abbreviation     string      `json:"abbreviation"`
	GmtOffset        int         `json:"gmtOffset"`
	Dst              string      `json:"dst"`
	ZoneStart        int         `json:"zoneStart"`
	ZoneEnd          interface{} `json:"zoneEnd"`
	NextAbbreviation interface{} `json:"nextAbbreviation"`
	Timestamp        int         `json:"timestamp"`
	Formatted        string      `json:"formatted"`
}

func AddCity(city ApiRequests.City) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into cities(cityName, Lat, Lon, WasUpdated, TimeZone) values(?, ?, ?, ?, ?)",
		city[0].Name, city[0].Lat, city[0].Lon, 0, GetTimeZone(city))
	if err != nil {
		panic(err)
	}

}
func AddUser(username string, chatid int) {
	if CheckIfUserExists(chatid) {
		return
	}
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into users(chatid, relocationsLeft, relocationsTries, hasSub) values(?, 5, 5, 0)",
		username, chatid)
	if err != nil {
		panic(err)
	}
}
func CheckIfUserExists(chatid int) bool {
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
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row := db.QueryRow("select exists(select * from cities where city = ?)", cityName)
	exist := false
	row.Scan(&exist)
	return exist
}

func TryToSetUserCity(cityName string, chatid int) string {
	if cityName == GetUserCityName(chatid) {
		return "Этот город уже указан"
	}
	if CheckIfCityExistsInDB(cityName) {
		SetUserCity(cityName, chatid)
		return "Учпечно"
	}
	if city, err := ApiRequests.CheckIfCityIsReal(cityName); err == false {
		AddCity(city)
		SetUserCity(cityName, chatid)
		return "Учпечно"
	}
	return "Кажется, этого города не существует..."
}

func SetUserCity(cityName string, chatid int) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if !CheckWeatherActuality(cityName) {
		_, err = db.Exec("update users set city = ? where chatid = ?", cityName, chatid)
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
	_, err = db.Exec("update cities set todaysWeather = ? where cityName = ?", ApiRequests.CheckTodaysWeather(lat, lon))
	if err != nil {
		panic(err)
	}
}
func SetWeatherActualityTrue(cityName string) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("update cities set wasUpdated = 1 where cityName = ?", cityName)
	if err != nil {
		panic(err)
	}
}
func SetWeatherActualityFalse() {
	//
}

// вызывается в server при попытке запроса погоды, затем, если город указан, апдейтит погоду и возвращает юзеру результат
func GetUserCityName(chatid int) string {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	row, err := db.Query("select cityName from users where chatid = ?", chatid)
	if err != nil {
		panic(err)
	}
	var cityName string
	for row.Next() {
		err := row.Scan(&cityName)
		if err != nil {
			panic(err)
		}
	}
	return cityName
}
func GetTimeZone(city ApiRequests.City) int {
	timeZoneTime := 0
	url := fmt.Sprintf("http://api.timezonedb.com/v2.1/get-time-zone?key=%v&format=json&by=position&lat=%v&lng=%v", ApiTokens.TimeZone, city[0].Lat, city[0].Lon)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Panic(err)
	}
	var timeZone TimeZone
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(body, &timeZone)
	timeZoneTime = timeZone.GmtOffset / 3600
	fmt.Println(timeZone)
	return timeZoneTime
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
8) Проверка, что время == XX:01
9) Обновление статуса на 0 для городов, у которых 12 ночи по местному времени(sql запрос)
9.1) Обновление relocationsLeft в начале месяца(не по местному)
9.2) Обновление relocationsTries в начале дня(не по местному)
10)+ Узнать местное время в новом(для БД) городе(api request)
11)+ Проверка на наличие пользователя в бд(при нажатии на /start)(sql запрос)
12)+ Добавление пользователя в БД(sql запрос)
13) Уменьшение relocationsLeft при успешной смене города
13.5) Уменьшение relocationsTries при не успешной смене города
*/
