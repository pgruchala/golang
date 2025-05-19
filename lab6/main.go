package main

import (
	"fmt"
	"lab6/data"
	_ "lab6/indicators/momentum"
	_ "lab6/indicators/trend"
	_ "lab6/indicators/volatility"
	"os"
)

func main() {
	fmt.Println("---------Analiza Danych Giełgowych---------")
	filePath := getFilePath()
	loader := data.NewCSVLoader()
	entries, err := loader.LoadData(filePath)
	if err != nil {
		fmt.Errorf("błąd wczytywanie danych %v", err)
		os.Exit(1)
	}
	fmt.Println("Wczytywanie zakończone pomyślnie")
	prices := extractPrices(entries)
	for {
		fmt.Println("\nWybierz opcję:")
		fmt.Println("1. Wskaźnik trendu: SMA (Simple Moving Average)")
		fmt.Println("2. Wskaźnik zmienności: Wstęgi Bollingera")
		fmt.Println("3. Wskaźnik impetu: CCI (Commodity Channel Index)")
		fmt.Println("4. Pokaż dane źródłowe")
		fmt.Println("5. Wyjście")
	
		choice := getUserChoice()
		switch choice {
		case 1:
			calculateSMA(entries, prices.closes)
		case 2:
			calculateBollinger(entries, prices.closes)
		case 3:
			calculateCCI(entries, prices.highs, prices.lows, prices.closes)
		case 4:
			displayCSV(entries)
		case 5:
			fmt.Println("Do widzenia!")
			return
		default:
			fmt.Println("Niepoprawny wybór")
		}
	}
}