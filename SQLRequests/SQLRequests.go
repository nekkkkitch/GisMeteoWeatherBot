package SQLRequests

import (
	GisMeteoRequest "WeatherBot/GisMeteoRequests"
	ApiTokens "WeatherBot/key"
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

}
func CheckForCity(cityName string) bool {
	return false
}
func AddUser() {

}
func CheckWeatherActuality() bool {
	return false
}
func UpdateWeather() {

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
