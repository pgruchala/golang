package main

import (
	"flag"
	"fmt"
	"log"
)

func handleCurrent(cmd *flag.FlagSet, args []string) {
	cmd.Parse(args)
	if cmd.NArg() == 0 {
		log.Fatal("Błąd: Nie podano nazwy miasta.")
	}
	cityName := cmd.Arg(0)
	city, found := findCityByName(cities, cityName)
	if !found {
		log.Fatalf("Błąd: Nie znaleziono miasta '%s'", cityName)
	}

	fmt.Printf("Pobieranie aktualnej pogody dla: %s\n", city.Name)
	weather, err := fetchCurrentWeather(city.Latitude, city.Longitude)
	if err != nil {
		log.Fatalf("Błąd pobierania pogody: %v", err)
	}
	displayCurrentWeather(*weather, city.Name)
}


func handleForecast(cmd *flag.FlagSet, args []string, days *int, config Config) {
	cmd.Parse(args)
	if cmd.NArg() == 0 {
		log.Fatal("Błąd: Nie podano nazwy miasta.")
	}
	cityName := cmd.Arg(0)
	city, found := findCityByName(cities,cityName)
	if !found {
		log.Fatalf("Błąd: Nie znaleziono miasta '%s'", cityName)
	}

	fmt.Printf("Pobieranie prognozy pogody na %d dni dla: %s\n", *days, city.Name)
	weather, err := fetchForecast(city.Latitude, city.Longitude, *days)
	if err != nil {
		log.Fatalf("Błąd pobierania prognozy: %v", err)
	}
	displayForecast(*weather, city.Name)
	checkForExtremeWeather(*weather, config)
	err = generatePlot(*weather, "prognoza.png")
	if err != nil {
		log.Fatalf("Błąd generowania wykresu: %v", err)
	}
	fmt.Println("\nWykres prognozy został zapisany do pliku 'prognoza.png'")
}
func handleHistory(cmd *flag.FlagSet, args []string, start, end *string) {
    cmd.Parse(args)

    
    if *start == "" || *end == "" {
        log.Fatal("Błąd: Wymagane są daty --start i --end (format RRRR-MM-DD).")
    }
    
    if cmd.NArg() == 0 {
        log.Fatal("Błąd: Nie podano nazwy miasta.")
    }
	cityName := cmd.Arg(0)
	city, found := findCityByName(cities,cityName)
	if !found {
		log.Fatalf("Błąd: Nie znaleziono miasta '%s'", cityName)
	}

	fmt.Printf("Pobieranie historycznej pogody dla: %s (od %s do %s)\n", city.Name, *start, *end)
	weather, err := fetchHistorical(city.Latitude, city.Longitude, *start, *end)
	if err != nil {
		log.Fatalf("Błąd pobierania danych historycznych: %v", err)
	}
	displayHistory(*weather, city.Name)
	err = generatePlot(*weather, "historia.png")
	if err != nil {
		log.Fatalf("Błąd generowania wykresu: %v", err)
	}
	fmt.Println("\nWykres danych historycznych został zapisany do pliku 'historia.png'")
}