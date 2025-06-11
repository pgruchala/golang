package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var cities []City

type City struct {
	Name      string
	Latitude  float64
	Longitude float64
}

func getWeatherDescription(code int) string {
	switch code {
	case 0:
		return "Czyste niebo"
	case 1, 2, 3:
		return "Głównie bezchmurnie, częściowe zachmurzenie"
	case 45, 48:
		return "Mgła"
	case 51, 53, 55:
		return "Mżawka"
	case 61, 63, 65:
		return "Deszcz"
	case 66, 67:
		return "Marznący deszcz"
	case 71, 73, 75:
		return "Opady śniegu"
	case 80, 81, 82:
		return "Przelotne opady deszczu"
	case 85, 86:
		return "Przelotne opady śniegu"
	case 95, 96, 99:
		return "Burza"
	default:
		return "Nieznany"
	}
}
func checkForExtremeWeather(data WeatherResponse, config Config) {
	fmt.Println("\nAnaliza zagrożeń pogodowych:")
	found := false
	for i, date := range data.Daily.Time {
		tempMax := data.Daily.TempMax[i]
		tempMin := data.Daily.TempMin[i]

		if tempMax > config.ExtremeTempMax {
			fmt.Printf("  - UWAGA [%s]: Przewidywana ekstremalnie wysoka temperatura: %.1f °C (próg: %.1f °C)\n", date, tempMax, config.ExtremeTempMax)
			found = true
		}
		if tempMin < config.ExtremeTempMin {
			fmt.Printf("  - UWAGA [%s]: Przewidywana ekstremalnie niska temperatura: %.1f °C (próg: %.1f °C)\n", date, tempMin, config.ExtremeTempMin)
			found = true
		}
	}
	if !found {
		fmt.Println("  - Brak przewidywanych ekstremalnych zjawisk pogodowych.")
	}
}

func findCityByName(cities []City, name string) (City, bool) {
	for _, city := range cities {
		if strings.EqualFold(city.Name, name) {
			return city, true
		}
	}
	return City{}, false
}

//loading config from .json

type Config struct {
	ExtremeTempMax float64 `json:"extreme_temp_max"`
	ExtremeTempMin float64 `json:"extreme_temp_min"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
