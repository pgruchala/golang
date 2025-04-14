package main

import (
	"fmt"
	"sync"
)

func main() {
	// Liczba zamówień do wygenerowania
	numOrders := 50
	
	// Liczba równoległych pracowników
	numWorkers := 5
	
	// Kanały do komunikacji
	ordersChan := make(chan Order, numOrders)
	resultsChan := make(chan ProcessResult, numOrders*2) // *2 bo mogą być ponowne próby
	failedOrdersChan := make(chan Order, numOrders)
	
	var wg sync.WaitGroup
	
	// Uruchamianie generatora zamówień
	go orderGenerator(numOrders, ordersChan)
	
	// Uruchamianie pracowników
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, ordersChan, resultsChan, &wg)
	}
	
	// Uruchamianie kolektora wyników
	wg.Add(1)
	go collectResults(resultsChan, failedOrdersChan, &wg, numOrders)
	
	// Uruchamianie ponownych prób dla nieudanych zamówień
	wg.Add(1)
	go retryFailedOrders(failedOrdersChan, resultsChan, &wg)
	
	// Czekanie na zakończenie wszystkich gorutyn
	wg.Wait()
	
	fmt.Println("\nSymulacja zakończona")
}