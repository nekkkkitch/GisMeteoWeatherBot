package GisMeteoRequest

import (
	ApiTokens "WeatherBot/key"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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
	Data struct {
		Astro struct {
			Sun struct {
				Sunrise time.Time `json:"sunrise"`
				Sunset  time.Time `json:"sunset"`
				Polar   any       `json:"polar"`
			} `json:"sun"`
			Moon struct {
				NextFull           time.Time `json:"next_full"`
				PreviousFull       time.Time `json:"previous_full"`
				Phase              string    `json:"phase"`
				PercentIlluminated float64   `json:"percent_illuminated"`
			} `json:"moon"`
		} `json:"astro"`
		Icon struct {
			IconWeather string `json:"icon-weather"`
			Emoji       string `json:"emoji"`
		} `json:"icon"`
		Kind        string `json:"kind"`
		Description string `json:"description"`
		Date        struct {
			Utc            time.Time `json:"UTC"`
			Local          time.Time `json:"local"`
			Unix           int       `json:"unix"`
			TimeZoneOffset int       `json:"timeZoneOffset"`
		} `json:"date"`
		City struct {
			Name      string  `json:"name"`
			NameP     string  `json:"nameP"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"city"`
		Wind struct {
			Direction struct {
				Degree int `json:"degree"`
				Scale8 int `json:"scale_8"`
			} `json:"direction"`
			Speed struct {
				MS float64 `json:"m_s"`
			} `json:"speed"`
			GustSpeed struct {
				MS float64 `json:"m_s"`
			} `json:"gust_speed"`
			AlternateDirection bool `json:"alternate_direction"`
		} `json:"wind"`
		Precipitation struct {
			Type      int `json:"type"`
			TypeExt   int `json:"type_ext"`
			Amount    int `json:"amount"`
			Intensity int `json:"intensity"`
			Duration  int `json:"duration"`
		} `json:"precipitation"`
		Temperature struct {
			Air struct {
				C float64 `json:"C"`
			} `json:"air"`
			Comfort struct {
				C float64 `json:"C"`
			} `json:"comfort"`
			Water struct {
				C float64 `json:"C"`
			} `json:"water"`
		} `json:"temperature"`
		Storm struct {
			Cape       float64 `json:"cape"`
			Prediction bool    `json:"prediction"`
		} `json:"storm"`
		Cloudiness struct {
			Percent int `json:"percent"`
			Scale3  int `json:"scale_3"`
		} `json:"cloudiness"`
		Visibility struct {
			Horizontal struct {
				M int `json:"m"`
			} `json:"horizontal"`
		} `json:"visibility"`
		Humidity struct {
			Percent  int `json:"percent"`
			DewPoint struct {
				C float64 `json:"C"`
			} `json:"dew_point"`
		} `json:"humidity"`
		Pressure struct {
			MmHgAtm int `json:"mm_hg_atm"`
		} `json:"pressure"`
	} `json:"data"`
	Jsonapi struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	Meta struct {
		Status     bool `json:"status"`
		StatusCode int  `json:"status_code"`
	} `json:"meta"`
}
type TodaysWeather struct {
	Data []struct {
		Astro struct {
			Sun struct {
				Sunrise time.Time `json:"sunrise"`
				Sunset  time.Time `json:"sunset"`
				Polar   any       `json:"polar"`
			} `json:"sun"`
			Moon struct {
				NextFull           time.Time `json:"next_full"`
				PreviousFull       time.Time `json:"previous_full"`
				Phase              string    `json:"phase"`
				PercentIlluminated float64   `json:"percent_illuminated"`
			} `json:"moon"`
		} `json:"astro"`
		Icon struct {
			IconWeather string `json:"icon-weather"`
			Emoji       string `json:"emoji"`
		} `json:"icon"`
		Kind        string `json:"kind"`
		Description string `json:"description"`
		Date        struct {
			Utc            time.Time `json:"UTC"`
			Local          time.Time `json:"local"`
			Unix           int       `json:"unix"`
			TimeZoneOffset int       `json:"timeZoneOffset"`
		} `json:"date"`
		City struct {
			Name      string  `json:"name"`
			NameP     string  `json:"nameP"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"city"`
		Wind struct {
			Direction struct {
				Degree int `json:"degree"`
				Scale8 int `json:"scale_8"`
			} `json:"direction"`
			Speed struct {
				MS float64 `json:"m_s"`
			} `json:"speed"`
			GustSpeed struct {
				MS float64 `json:"m_s"`
			} `json:"gust_speed"`
			AlternateDirection bool `json:"alternate_direction"`
		} `json:"wind"`
		Precipitation struct {
			Type      int     `json:"type"`
			TypeExt   int     `json:"type_ext"`
			Amount    float64 `json:"amount"`
			Intensity int     `json:"intensity"`
			Duration  int     `json:"duration"`
		} `json:"precipitation"`
		Temperature struct {
			Air struct {
				C float64 `json:"C"`
			} `json:"air"`
			Comfort struct {
				C float64 `json:"C"`
			} `json:"comfort"`
			Water struct {
				C float64 `json:"C"`
			} `json:"water"`
		} `json:"temperature"`
		Storm struct {
			Cape       float64 `json:"cape"`
			Prediction bool    `json:"prediction"`
		} `json:"storm"`
		Cloudiness struct {
			Percent int `json:"percent"`
			Scale3  int `json:"scale_3"`
		} `json:"cloudiness"`
		Visibility struct {
			Horizontal struct {
				M int `json:"m"`
			} `json:"horizontal"`
		} `json:"visibility"`
		Humidity struct {
			Percent  int `json:"percent"`
			DewPoint struct {
				C float64 `json:"C"`
			} `json:"dew_point"`
		} `json:"humidity"`
		Pressure struct {
			MmHgAtm int `json:"mm_hg_atm"`
		} `json:"pressure"`
	} `json:"data"`
	Jsonapi struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	Meta struct {
		Status     bool `json:"status"`
		StatusCode int  `json:"status_code"`
	} `json:"meta"`
}

func CheckCurrentWeather(city City) (CurrentWeather, string) {
	var weather CurrentWeather
	if len(city) == 0 {
		return weather, "Город с таким названием не найден, попробуйте ввести другой."
	}
	weatherUrl := fmt.Sprintf("https://api.gismeteo.net/v3/weather/current/?latitude=%v&longitude=%v", city[0].Lat, city[0].Lon) //получение погоды на сейчас
	weatherReq, err := http.NewRequest("GET", weatherUrl, nil)
	weatherReq.Header.Add("X-Gismeteo-Token", ApiTokens.GisMeteoToken)
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
	weatherUrl := fmt.Sprintf("https://api.gismeteo.net/v3/weather/forecast/h3/?latitude=%v&longitude=%v", city[0].Lat, city[0].Lon) //получение погоды на сейчас
	weatherReq, err := http.NewRequest("GET", weatherUrl, nil)
	weatherReq.Header.Add("X-Gismeteo-Token", ApiTokens.GisMeteoToken)
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
