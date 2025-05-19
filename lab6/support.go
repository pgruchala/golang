package main

import (
	"bufio"
	"fmt"
	"lab6/data"
	"lab6/indicators/momentum"
	"lab6/indicators/trend"
	"lab6/indicators/volatility"
	"math"
	"os"
	"strconv"
	"strings"
)

type PriceData struct {
	closes []float64
	opens  []float64
	highs  []float64
	lows   []float64
}

func getFilePath() string {
	reader := bufio.NewReader(os.Stdin)
	var filePath string
	for {
		fmt.Print("Podaj ścieżkę do pliku z danymi")
		filePath, _ = reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)

		if _, err := os.Stat(filePath); err == nil {
			break
		} else {
			fmt.Println("Nie można odnaleźć pliku. Spróbuj ponownie.")
		}
	}
	return filePath
}
func extractPrices(entries []data.Entry) PriceData {
	result := PriceData{
		closes: make([]float64, len(entries)),
		opens:  make([]float64, len(entries)),
		highs:  make([]float64, len(entries)),
		lows:   make([]float64, len(entries)),
	}
	for i, entry := range entries {
		result.closes[i] = entry.Close
		result.opens[i] = entry.Open
		result.highs[i] = entry.High
		result.lows[i] = entry.Low
	}
	return result
}

func getUserChoice() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Wybór:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err == nil {
			return choice
		}
		fmt.Println("NIeprawidłowy wybór")
	}
}
func getPositiveInt(prompt string) int {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		value, err := strconv.Atoi(input)
		if err == nil && value > 0 {
			return value
		}

		fmt.Println("Nieprawidłowa wartość. Podaj dodatnią liczbę całkowitą.")
	}
}

func getPositiveFloat(prompt string) float64 {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		value, err := strconv.ParseFloat(input, 64)
		if err == nil && value > 0 {
			return value
		}

		fmt.Println("Nieprawidłowa wartość. Podaj dodatnią liczbę.")
	}
}

func calculateSMA(entries []data.Entry, prices []float64) {
	period := getPositiveInt("Podaj okres dla SMA: ")

	smaValues, err := trend.CalculateSMA(prices, period)
	if err != nil {
		fmt.Printf("Błąd podczas obliczania SMA: %v\n", err)
		return
	}

	fmt.Printf("\nWyniki SMA (okres %d):\n", period)
	fmt.Println("Data\t\tCena zamknięcia\tSMA")
	fmt.Println("------------------------------------------")

	offset := len(entries) - len(smaValues)

	for i := 0; i < len(smaValues); i++ {
		entry := entries[i+offset]
		fmt.Printf("%s\t%.2f\t\t%.2f\n",
			entry.Date.Format("2006-01-02"),
			entry.Close,
			smaValues[i])
	}
}

func calculateBollinger(entries []data.Entry, prices []float64) {
	period := getPositiveInt("Podaj okres dla Wstęg Bollingera: ")
	stdDev := getPositiveFloat("Podaj mnożnik odchylenia standardowego (typowo 2.0): ")

	bbValues, err := volatility.CalculateBollinger(prices, period, stdDev)
	if err != nil {
		fmt.Printf("Błąd podczas obliczania Wstęg Bollingera: %v\n", err)
		return
	}

	fmt.Printf("\nWyniki Wstęg Bollingera (okres %d, odchylenie %.2f):\n", period, stdDev)
	fmt.Println("Data\t\tCena\t\tGórna\t\tŚrodkowa\tDolna")
	fmt.Println("-------------------------------------------------------------")

	for i := period - 1; i < len(entries); i++ {
		entry := entries[i]
		bb := bbValues[i]
		if !math.IsNaN(bb.MiddleBand) && !math.IsNaN(bb.UpperBand) && !math.IsNaN(bb.LowerBand) {
			fmt.Printf("%s\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\n",
				entry.Date.Format("2006-01-02"),
				entry.Close,
				bb.UpperBand,
				bb.MiddleBand,
				bb.LowerBand)
		}
	}
}

func calculateCCI(entries []data.Entry, highs, lows, closes []float64) {
	period := getPositiveInt("Podaj okres dla CCI: ")

	cciValues, err := momentum.CalculateCCI(highs, lows, closes, period)
	if err != nil {
		fmt.Printf("Błąd podczas obliczania CCI: %v\n", err)
		return
	}

	fmt.Printf("\nWyniki CCI (okres %d):\n", period)
	fmt.Println("Data\t\tCena zamknięcia\tCCI")
	fmt.Println("------------------------------------------")

	offset := len(entries) - len(cciValues)

	for i := 0; i < len(cciValues); i++ {
		entry := entries[i+offset]
		fmt.Printf("%s\t%.2f\t\t%.2f\n",
			entry.Date.Format("2006-01-02"),
			entry.Close,
			cciValues[i])
	}
}

func displayCSV(entries []data.Entry) {
	fmt.Println("\nDane źródłowe:")
	fmt.Println("Data\t\tOtwarcie\tNajwyższa\tNajniższa\tZamknięcie\tWolumen")
	fmt.Println("-------------------------------------------------------------------------------")

	startIdx := 0
	if len(entries) > 20 {
		startIdx = len(entries) - 20
		fmt.Println("(Wyświetlono tylko 20 ostatnich wpisów...)")
	}

	for i := startIdx; i < len(entries); i++ {
		e := entries[i]
		fmt.Printf("%s\t%.2f\t\t%.2f\t\t%.2f\t\t%.2f\t\t%.0f\n",
			e.Date.Format("2006-01-02"),
			e.Open,
			e.High,
			e.Low,
			e.Close,
			e.Volume)
	}
}
