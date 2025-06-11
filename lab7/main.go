package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	cities = loadData()
	if cities == nil {
		log.Fatal("NIe udało się załadować informacji o miastach")
	}
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("Błąd wczytywania konfiguracji", err)
	}
	switch os.Args[1] {
	case "aktualna":
		currentCmd := flag.NewFlagSet("aktualna", flag.ExitOnError)
		handleCurrent(currentCmd, os.Args[2:])
	case "prognoza":
		forecastCmd := flag.NewFlagSet("prognoza", flag.ExitOnError)
		forecastDays := forecastCmd.Int("dni", 7, "Liczba dni prognozy")
		handleForecast(forecastCmd, os.Args[2:], forecastDays, config)
	case "historia":
		historyCmd := flag.NewFlagSet("historia", flag.ExitOnError)
		historyStartDate := historyCmd.String("start", "", "Data początkowa (RRRR-MM-DD)")
		historyEndDate := historyCmd.String("end", "", "Data końcowa (RRRR-MM-DD)")
		handleHistory(historyCmd, os.Args[2:], historyStartDate, historyEndDate)
	default:
		printHelp()
	}

}
func printHelp() {
	fmt.Println("Użycie: pogoda <komenda> [opcje] <miasto>")
	fmt.Println("\nDostępne komendy:")
	fmt.Println("  aktualna <miasto>                - Aktualne dane pogodowe")
	fmt.Println("  prognoza <miasto> --dni <liczba> - Prognoza pogody (domyślnie 7 dni)")
	fmt.Println("  historia <miasto> --daty <daty>  - Historyczne dane (daty: RRRR-MM-DD,RRRR-MM-DD)")
}
