package SQLRequests

import (
	GisMeteoRequest "WeatherBot/GisMeteoRequests"
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
	ID               int                  `json:"ID"`
	Username         string               `json:"Username"`
	ChatID           int                  `json:"ChatID"`
	City             GisMeteoRequest.City `json:"City"`
	TimeNotification time.Time            `json:"TimeNotification"`
	CallsLimit       int                  `json:"CallsLimit"`
	CallsLeft        int                  `json:"CallsLeft"`
	HasSub           bool                 `json:"HasSub"`
}
type City struct {
	CityName       string                         `json:"CityName"`
	WeatherToday   GisMeteoRequest.TodaysWeather  `json:"WeatherToday"`
	WeatherCurrent GisMeteoRequest.CurrentWeather `json:"WeatherCurrent"`
	Lat            float64                        `json:"Lat"`
	Lon            float64                        `json:"Lon"`
	WasUpdated     bool                           `json:"WasUpdated"`
	TimeZone       int                            `json:"TimeZone"`
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

func AddCity(city GisMeteoRequest.City) {
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into cities(cityName, Lat, Lon, WasUpdated, TimeZone) values(?, ?, ?, ?, ?)",
		city[0].Name, city[0].Lat, city[0].Lon, GetTimeZone(city), 1)
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
	_, err = db.Exec("insert into users(username, chatid, relocationsLeft, hasSub) values(?, ?, 5, 0)",
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
func CheckIfCityExists(cityName string) bool {
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

// проверяет существование города(возвращает ошибку, если нет) -> проверяет существование города в базе(добавляет, если нет)-> устанавливает город для юзера
func SetUserCity(chatID int, cityName string) string {
	city, problem := GisMeteoRequest.CheckIfCityIsReal(cityName) // проверка что такой город существует
	if problem != "" {
		return problem
	}
	if !CheckIfCityExists(cityName) { // проверка на наличие города в базе
		AddCity(city)
	}
	db, err := sql.Open("mysql", ApiTokens.SQLOpening)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("update users set city = ? where chatid = ?", cityName, chatID)
	if err != nil {
		panic(err.Error())
	}
	return ""
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
func UpdateWeather(chatid int) { //todo: dodelat
	cityName, problem := GetUserCityName(chatid)

}

// вызывается в server при попытке запроса погоды, затем, если город указан, апдейтит погоду и возвращает юзеру результат
// коммент говно надо переосмыслить
func GetUserCityName(chatid int) (string, string) {
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
	if cityName == "" {
		return "", "Кажется, вы забыли указать свой город. Сделайте это по кнопке ниже, пожалуйста"
	}
	return cityName, ""
}
func GetTimeZone(city GisMeteoRequest.City) int {
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

//todo: обнулять переменную city/wasupdated в 00:00 по местному времени
//если пользователь запрашивает данные с city где wasupdate == false(0), то запрашивать нынешнюю погоду в данном городе
//может быть запретить менять город чаще n раз за месяц
