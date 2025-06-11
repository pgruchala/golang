package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.open-meteo.com/v1/forecast"

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Temperature         float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		WindSpeed           float64 `json:"wind_speed_10m"`
		WeatherCode         int     `json:"weather_code"`
		UvIndex             float64 `json:"uv_index"`
		IsDay               int     `json:"is_day"`
	} `json:"current"`
	Daily struct {
		Time        []string  `json:"time"`
		WeatherCode []int     `json:"weathercode"`
		TempMax     []float64 `json:"temperature_2m_max"`
		TempMin     []float64 `json:"temperature_2m_min"`
		Sunrise     []string  `json:"sunrise"`
		Sunset      []string  `json:"sunset"`
	} `json:"daily"`
}

func fetchCurrentWeather(lat, lon float64) (*WeatherResponse, error) {
	params := "current=temperature_2m,apparent_temperature,wind_speed_10m,weather_code,uv_index,is_day"
	url := fmt.Sprintf("%s?latitude=%.4f&longitude=%.4f&%s", baseURL, lat, lon, params)
	return fetchWeatherData(url)
}

func fetchForecast(lat, lon float64, days int) (*WeatherResponse, error) {
	params := "daily=weathercode,temperature_2m_max,temperature_2m_min,sunrise,sunset"
	url := fmt.Sprintf("%s?latitude=%.4f&longitude=%.4f&%s&forecast_days=%d", baseURL, lat, lon, params, days)
	return fetchWeatherData(url)
}

func fetchHistorical(lat, lon float64, start, end string) (*WeatherResponse, error) {

	_, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, fmt.Errorf("nieprawidłowy format daty początkowej: %s", start)
	}
	_, err = time.Parse("2006-01-02", end)
	if err != nil {
		return nil, fmt.Errorf("nieprawidłowy format daty końcowej: %s", end)
	}

	params := "daily=weathercode,temperature_2m_max,temperature_2m_min,sunrise,sunset"
	url := fmt.Sprintf("%s?latitude=%.4f&longitude=%.4f&%s&start_date=%s&end_date=%s", baseURL, lat, lon, params, start, end)
	return fetchWeatherData(url)
}

func fetchWeatherData(url string) (*WeatherResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("problem z połączeniem z API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("błąd API, status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu odpowiedzi API: %w", err)
	}

	var weatherData WeatherResponse
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("błąd parsowania danych JSON: %w", err)
	}

	return &weatherData, nil
}
