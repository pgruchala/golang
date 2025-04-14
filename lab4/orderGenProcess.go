package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
	"errors"
)

const (	
	numOrders        = 20 
	numWorkers       = 5  
	maxRetries       = 2  
	baseSuccessRate  = 0.8
	retrySuccessRate = 0.6
)

func orderGenerator(ordersChan chan<- Order, wg *sync.WaitGroup) {
	defer wg.Done() 
	defer close(ordersChan)

	fmt.Println("Generator zamówień: Start")
	for i := 1; i <= numOrders; i++ {
		order := generateNewOrder(i)
		fmt.Printf("Generator: Wygenerowano zamówienie ID %d dla %s (%.2f PLN)\n", order.ID, order.CustomerName, order.TotalAmount)
		ordersChan <- order

		
		sleepDuration := time.Duration(rand.IntN(500)+50) * time.Millisecond
		time.Sleep(sleepDuration)
	}
	fmt.Println("Generator zamówień: Zakończono generowanie.")
}

func processOrderWorker(id int, ordersChan <-chan Order, resultsChan chan<- ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done() // Zasygnalizuj zakończenie pracy workera

	fmt.Printf("Worker %d: Start\n", id)
	for order := range ordersChan { // Pętla działa dopóki kanał ordersChan nie zostanie zamknięty i opróżniony
		fmt.Printf("Worker %d: Przetwarzanie zamówienia ID %d dla %s\n", id, order.ID, order.CustomerName)

		var result ProcessResult
		success := false
		var processingErr error
		var totalProcessingTime time.Duration
		attempts := 0

		for attempts = 1; attempts <= maxRetries+1; attempts++ {
			// startTime := time.Now()
			processingDuration := time.Duration(rand.IntN(500)+100) * time.Millisecond
			time.Sleep(processingDuration)
			totalProcessingTime += processingDuration

			successRate := baseSuccessRate
			if attempts > 1 {
				successRate = retrySuccessRate 
				fmt.Printf("Worker %d: Ponowna próba (%d/%d) zamówienia ID %d\n", id, attempts-1, maxRetries, order.ID)
			}

			if rand.Float64() < successRate {
				success = true
				processingErr = nil 
				fmt.Printf("Worker %d: Sukces przetwarzania zamówienia ID %d (próba %d)\n", id, order.ID, attempts)
				break 
			} else {
				success = false
				processingErr = errors.New("błąd symulacji przetwarzania") 
				fmt.Printf("Worker %d: Niepowodzenie przetwarzania zamówienia ID %d (próba %d)\n", id, order.ID, attempts)
				if attempts > maxRetries {
					fmt.Printf("Worker %d: Osiągnięto limit prób dla zamówienia ID %d\n", id, order.ID)
					break
				}
				
				time.Sleep(time.Duration(rand.IntN(100)+50) * time.Millisecond)
			}
		}

		result = ProcessResult{
			OrderID:      order.ID,
			CustomerName: order.CustomerName,
			Success:      success,
			ProcessTime:  totalProcessingTime, 
			Error:        processingErr,       
		}

		resultsChan <- result
	}
	fmt.Printf("Worker %d: Koniec pracy (kanał zamówień zamknięty)\n", id)
}

func resultsCollector(resultsChan <-chan ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done() // Zasygnalizuj zakończenie pracy kolektora

	fmt.Println("Kolektor wyników: Start")
	var successCount, failureCount int
	var totalProcessed int

	for result := range resultsChan { // Pętla działa dopóki kanał resultsChan nie zostanie zamknięty i opróżniony
		totalProcessed++
		fmt.Printf("Kolektor: Otrzymano wynik dla zamówienia ID %d (Klient: %s, Sukces: %t, Czas: %v, Błąd: %v)\n",
			result.OrderID, result.CustomerName, result.Success, result.ProcessTime, result.Error)
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	// Obliczanie statystyk
	fmt.Println("\n--- Statystyki końcowe ---")
	fmt.Printf("Całkowita liczba przetworzonych (lub prób przetworzenia) zamówień: %d\n", totalProcessed)
	fmt.Printf("Zamówienia przetworzone pomyślnie: %d\n", successCount)
	fmt.Printf("Zamówienia zakończone niepowodzeniem (po %d próbach): %d\n", maxRetries, failureCount)

	if totalProcessed > 0 {
		successRate := float64(successCount) / float64(totalProcessed) * 100
		failureRate := float64(failureCount) / float64(totalProcessed) * 100
		fmt.Printf("Procent zamówień zakończonych sukcesem: %.2f%%\n", successRate)
		fmt.Printf("Procent zamówień zakończonych niepowodzeniem: %.2f%%\n", failureRate)
	} else {
		fmt.Println("Nie przetworzono żadnych zamówień.")
	}
	fmt.Println("Kolektor wyników: Zakończono zbieranie.")
}
