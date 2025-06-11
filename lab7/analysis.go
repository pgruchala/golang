package main

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

func displayCurrentWeather(data WeatherResponse, cityName string) {
	table := tablewriter.NewWriter(os.Stdout)

	var dayNight string
	if data.Current.IsDay == 1 {
		dayNight = "Dzień"
	} else {
		dayNight = "Noc"
	}

	table.Header([]string{"Parametr", "Wartość"})
	table.Append([]string{"Pora dnia", dayNight})
	table.Append([]string{"Temperatura", fmt.Sprintf("%.1f °C", data.Current.Temperature)})
	table.Append([]string{"Temperatura odczuwalna", fmt.Sprintf("%.1f °C", data.Current.ApparentTemperature)})
	table.Append([]string{"Prędkość wiatru", fmt.Sprintf("%.1f km/h", data.Current.WindSpeed)})
	table.Append([]string{"Indeks UV", fmt.Sprintf("%.1f", data.Current.UvIndex)})
	table.Append([]string{"Warunki", getWeatherDescription(data.Current.WeatherCode)})

	table.Render()
}

func displayForecast(data WeatherResponse, cityName string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Data", "Temp. Max (°C)", "Temp. Min (°C)", "Wschód słońca", "Zachód słońca", "Opis"})

	for i, day := range data.Daily.Time {
		sunrise, _ := time.Parse("2006-01-02T15:04", data.Daily.Sunrise[i])
		sunset, _ := time.Parse("2006-01-02T15:04", data.Daily.Sunset[i])
		row := []string{
			day,
			fmt.Sprintf("%.1f", data.Daily.TempMax[i]),
			fmt.Sprintf("%.1f", data.Daily.TempMin[i]),
			sunrise.Format("15:04"),
			sunset.Format("15:04"),
			getWeatherDescription(data.Daily.WeatherCode[i]),
		}
		table.Append(row)
	}
	table.Render()
}
func displayHistory(data WeatherResponse, cityName string) {
	displayForecast(data, cityName) // Używamy tej samej funkcji co dla prognozy, bo format jest identyczny
}
