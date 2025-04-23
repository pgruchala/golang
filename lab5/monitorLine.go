package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func monitorTravelTime(stops []Stop, lineID int, stop1, stop2 *Stop) {
	fmt.Printf("\nMonitorowanie czasu przejazdu linii %d między przystankami:\n - %s\n - %s\n\n", lineID, stop1.StopName, stop2.StopName)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		deps1, err1 := fetchEstimatedDepartures(fmt.Sprintf("%d", stop1.StopID))
		deps2, err2 := fetchEstimatedDepartures(fmt.Sprintf("%d", stop2.StopID))

		if err1 != nil || err2 != nil {
			fmt.Println("Błąd podczas pobierania odjazdów:", err1, err2)
			return
		}
		var dep1, dep2 *EstimatedDeparture
		for _, d := range deps1 {
			if d.RouteID == lineID {
				dep1 = &d
				break
			}
		}
		for _, d := range deps2 {
			if d.RouteID == lineID {
				dep2 = &d
				break
			}
		}

		if dep1 == nil || dep2 == nil {
			fmt.Printf("Brak odjazdów linii %d z jednego z przystanków.\n", lineID)
		} else {
			t1, err1 := time.Parse(time.RFC3339, dep1.EstimatedTime)
			t2, err2 := time.Parse(time.RFC3339, dep2.EstimatedTime)

			if err1 != nil || err2 != nil {
				fmt.Println("Błąd parsowania czasu odjazdu.")
			} else {
				diff := t2.Sub(t1)
				if diff < 0 {
					fmt.Println("Uwaga: czas odjazdu z drugiego przystanku jest wcześniejszy niż z pierwszego.")
				} else {
					fmt.Printf("Szacowany czas przejazdu: %v minut %v sekund\n", int(diff.Minutes()), int(diff.Seconds())%60)
				}
			}
		}

		fmt.Println("Aktualizacja za 10 sekund...")
		<-ticker.C
	}
}

func getLineAndStops(stops []Stop) (int, *Stop, *Stop, error) {
	stopName1 := getUserInput("Podaj nazwę pierwszego przystanku:")
	stop1 := findStopByName(stops, stopName1)
	if stop1 == nil {
		return 0, nil, nil, fmt.Errorf("nie znaleziono pierwszego przystanku")
	}

	stopName2 := getUserInput("Podaj nazwę drugiego przystanku:")
	stop2 := findStopByName(stops, stopName2)
	if stop2 == nil {
		return 0, nil, nil, fmt.Errorf("nie znaleziono drugiego przystanku")
	}

	lineIDStr := getUserInput("Podaj numer linii:")
	lineID, err := strconv.Atoi(lineIDStr)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("nieprawidłowy numer linii")
	}
	return lineID, stop1, stop2, nil
}

func monitorTwoLines(stops []Stop) {
	fmt.Println("\n=== Informacje dla pierwszej linii ===")
	lineID1, stop1A, stop1B, err1 := getLineAndStops(stops)
	if err1 != nil {
		fmt.Println("Błąd dla pierwszej linii:", err1)
		return
	}

	fmt.Println("\n=== Informacje dla drugiej linii ===")
	lineID2, stop2A, stop2B, err2 := getLineAndStops(stops)
	if err2 != nil {
		fmt.Println("Błąd dla drugiej linii:", err2)
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("\n--- Rozpoczynam monitoring pierwszej linii ---")
		monitorTravelTime(stops, lineID1, stop1A, stop1B)
	}()

	go func() {
		defer wg.Done()
		fmt.Println("\n--- Rozpoczynam monitoring drugiej linii ---")
		monitorTravelTime(stops, lineID2, stop2A, stop2B)
	}()

	wg.Wait()
}