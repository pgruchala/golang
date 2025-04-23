package main

import (
	"fmt"
)

func main() {
	stopsURL := "https://ckan.multimediagdansk.pl/dataset/c24aa637-3619-4dc2-a171-a23eec8f2172/resource/d3e96eb6-25ad-4d6c-8651-b1eb39155945/download/stopsingdansk.json"

	stops, err := fetchStopsData(stopsURL)
	if err != nil {
		fmt.Println("Error fetching stops data:", err)
		return
	}
	for {
		fmt.Println("\nWybierz opcję:")
		fmt.Println("1 - Wyświetl informacje o przystanku")
		fmt.Println("2 - Monitoruj czas przejazdu między przystankami")
		fmt.Println("3 - Monitoruj czas przejazdu między przystankami dla dwóch jednoczesnie")
		fmt.Println("4 - Wyjście")

		opcja := getUserInput("Podaj numer opcji:")

		switch opcja {
		case "1":
			getStopInfo(stops)
		case "2":
			lineID, stop1, stop2, err := getLineAndStops(stops)
			if err != nil {
				fmt.Println("Błąd:", err)
				continue
			}
			monitorTravelTime(stops, lineID, stop1, stop2)
		case "3":
			monitorTwoLines(stops)
		case "4":
			fmt.Println("Koniec programu.")
			return
		default:
			fmt.Println("Nieprawidłowa opcja.")
		}
	}
}
