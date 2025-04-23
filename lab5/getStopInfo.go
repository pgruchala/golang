package main

import "fmt"

func getStopInfo(stops []Stop){
	przystanek := getUserInput("Podaj nazwę przystanku:")
	stop := findStopByName(stops, przystanek) 

	departures, err := fetchEstimatedDepartures(fmt.Sprintf("%d", stop.StopID))
	if err != nil {
		fmt.Println("Błąd podczas pobierania danych o odjazdach:", err)
		return
	}

	if len(departures) == 0 {
		fmt.Println("Brak estymowanych odjazdów dla tego przystanku.")
		return
	}

	fmt.Println("Nadchodzące odjazdy:")
	for _, dep := range departures {
		estimated := formatTime(dep.EstimatedTime)
		delay := dep.DelayInSeconds
		delayStr := ""
		if delay > 0 {
			delayStr = fmt.Sprintf(" (opóźnienie: %d sek)", delay)
		}
		fmt.Printf("Linia %d -> %s, odjazd o %s%s\n", dep.RouteID, dep.Headsign, estimated, delayStr)
	}
	
}