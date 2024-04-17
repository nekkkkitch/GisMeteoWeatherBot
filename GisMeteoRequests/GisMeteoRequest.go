package GisMeteoRequest

import (
	ApiTokens "WeatherBot/key"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type City []struct {
	Name       string `json:"name"`
	LocalNames struct {
		De          string `json:"de"`
		Et          string `json:"et"`
		Ko          string `json:"ko"`
		Sk          string `json:"sk"`
		Cs          string `json:"cs"`
		Pl          string `json:"pl"`
		Hu          string `json:"hu"`
		Ar          string `json:"ar"`
		Lt          string `json:"lt"`
		Sl          string `json:"sl"`
		Ca          string `json:"ca"`
		Uk          string `json:"uk"`
		Fi          string `json:"fi"`
		En          string `json:"en"`
		Ja          string `json:"ja"`
		Pt          string `json:"pt"`
		FeatureName string `json:"feature_name"`
		Ru          string `json:"ru"`
		Ro          string `json:"ro"`
		Es          string `json:"es"`
		ASCII       string `json:"ascii"`
		Fr          string `json:"fr"`
		Hr          string `json:"hr"`
	} `json:"local_names"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}
type CurrentWeather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}
type TodaysWeather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date      string `json:"date"`
			DateEpoch int    `json:"date_epoch"`
			Day       struct {
				MaxtempC          float64 `json:"maxtemp_c"`
				MaxtempF          float64 `json:"maxtemp_f"`
				MintempC          float64 `json:"mintemp_c"`
				MintempF          float64 `json:"mintemp_f"`
				AvgtempC          float64 `json:"avgtemp_c"`
				AvgtempF          float64 `json:"avgtemp_f"`
				MaxwindMph        float64 `json:"maxwind_mph"`
				MaxwindKph        float64 `json:"maxwind_kph"`
				TotalprecipMm     float64 `json:"totalprecip_mm"`
				TotalprecipIn     float64 `json:"totalprecip_in"`
				TotalsnowCm       float64 `json:"totalsnow_cm"`
				AvgvisKm          float64 `json:"avgvis_km"`
				AvgvisMiles       float64 `json:"avgvis_miles"`
				Avghumidity       int     `json:"avghumidity"`
				DailyWillItRain   int     `json:"daily_will_it_rain"`
				DailyChanceOfRain int     `json:"daily_chance_of_rain"`
				DailyWillItSnow   int     `json:"daily_will_it_snow"`
				DailyChanceOfSnow int     `json:"daily_chance_of_snow"`
				Condition         struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int    `json:"code"`
				} `json:"condition"`
				Uv float64 `json:"uv"`
			} `json:"day"`
			Astro struct {
				Sunrise          string `json:"sunrise"`
				Sunset           string `json:"sunset"`
				Moonrise         string `json:"moonrise"`
				Moonset          string `json:"moonset"`
				MoonPhase        string `json:"moon_phase"`
				MoonIllumination int    `json:"moon_illumination"`
				IsMoonUp         int    `json:"is_moon_up"`
				IsSunUp          int    `json:"is_sun_up"`
			} `json:"astro"`
			Hour []struct {
				TimeEpoch int     `json:"time_epoch"`
				Time      string  `json:"time"`
				TempC     float64 `json:"temp_c"`
				TempF     float64 `json:"temp_f"`
				IsDay     int     `json:"is_day"`
				Condition struct {
					Text string `json:"text"`
					Icon string `json:"icon"`
					Code int    `json:"code"`
				} `json:"condition"`
				WindMph      float64 `json:"wind_mph"`
				WindKph      float64 `json:"wind_kph"`
				WindDegree   int     `json:"wind_degree"`
				WindDir      string  `json:"wind_dir"`
				PressureMb   float64 `json:"pressure_mb"`
				PressureIn   float64 `json:"pressure_in"`
				PrecipMm     float64 `json:"precip_mm"`
				PrecipIn     float64 `json:"precip_in"`
				SnowCm       float64 `json:"snow_cm"`
				Humidity     int     `json:"humidity"`
				Cloud        int     `json:"cloud"`
				FeelslikeC   float64 `json:"feelslike_c"`
				FeelslikeF   float64 `json:"feelslike_f"`
				WindchillC   float64 `json:"windchill_c"`
				WindchillF   float64 `json:"windchill_f"`
				HeatindexC   float64 `json:"heatindex_c"`
				HeatindexF   float64 `json:"heatindex_f"`
				DewpointC    float64 `json:"dewpoint_c"`
				DewpointF    float64 `json:"dewpoint_f"`
				WillItRain   int     `json:"will_it_rain"`
				ChanceOfRain int     `json:"chance_of_rain"`
				WillItSnow   int     `json:"will_it_snow"`
				ChanceOfSnow int     `json:"chance_of_snow"`
				VisKm        float64 `json:"vis_km"`
				VisMiles     float64 `json:"vis_miles"`
				GustMph      float64 `json:"gust_mph"`
				GustKph      float64 `json:"gust_kph"`
				Uv           float64 `json:"uv"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func CheckCurrentWeather(city City) (CurrentWeather, string) {
	var weather CurrentWeather
	if len(city) == 0 {
		return weather, "Город с таким названием не найден, попробуйте ввести другой."
	}
	weatherUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%v&q=%v,%v&aqi=no", ApiTokens.WeatherToken, city[0].Lat, city[0].Lon) //получение погоды на сейчас
	weatherReq, err := http.NewRequest("GET", weatherUrl, nil)
	CheckForError(err)
	weatherResponse, err := http.DefaultClient.Do(weatherReq)
	CheckForError(err)
	weatherBody, err := io.ReadAll(weatherResponse.Body)
	CheckForError(err)
	err = json.Unmarshal(weatherBody, &weather)
	CheckForError(err)
	return weather, ""
}
func CheckTodaysWeather(city City) (TodaysWeather, string) {
	var weather TodaysWeather
	if len(city) == 0 {
		return weather, "Город с таким названием не найден, попробуйте ввести другой."
	}
	weatherUrl := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%v&q=%v,%v&aqi=no", ApiTokens.WeatherToken, city[0].Lat, city[0].Lon) //получение погоды на сейчас
	weatherReq, err := http.NewRequest("GET", weatherUrl, nil)
	CheckForError(err)
	weatherResponse, err := http.DefaultClient.Do(weatherReq)
	CheckForError(err)
	weatherBody, err := io.ReadAll(weatherResponse.Body)
	CheckForError(err)
	err = json.Unmarshal(weatherBody, &weather)
	CheckForError(err)
	return weather, ""
}
func UpdateCity(cityName string) (City, string) {
	var city City
	cityLocationUrl := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%v,643&appid=%v", cityName, ApiTokens.CityCoordsToken) // получение координат по городу
	cityReq, err := http.NewRequest("GET", cityLocationUrl, nil)
	if err != nil {
		return city, "Такого города нет в базе, проверьте правильность ввода или введите другой город."
	}
	cityResponse, err := http.DefaultClient.Do(cityReq)
	CheckForError(err)
	cityBody, err := io.ReadAll(cityResponse.Body)
	CheckForError(err)
	err = json.Unmarshal(cityBody, &city)
	CheckForError(err)
	fmt.Println(city)
	return city, ""
}
func CheckForError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
