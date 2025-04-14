package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

const (
	numOrders        = 20
	numWorkers       = 5
	maxRetries       = 2
	baseSuccessRate  = 0.8
	retrySuccessRate = 0.6
)

func orderGenerator(numOrders int, ordersChan chan<- Order) {
	for i := 1; i <= numOrders; i++ {
		order := generateNewOrder(i)
		fmt.Printf("Wygenerowano zamówienie #%d dla %s na kwotę %.2f zł\n", order.ID, order.CustomerName, order.TotalAmount)

		time.Sleep(time.Duration(rand.IntN(400)+100) * time.Millisecond)

		ordersChan <- order
	}
	close(ordersChan)
}

func processOrder(order Order) ProcessResult {
	processingTime := time.Duration(rand.IntN(700)+300) * time.Millisecond
	time.Sleep(processingTime)

	var err error
	success := true

	if rand.IntN(10) == 0 { //10% szans na błąd
		err = errors.New("błąd przetwarzania zamówienia")
		success = false
	}

	return ProcessResult{
		OrderID:      order.ID,
		CustomerName: order.CustomerName,
		Success:      success,
		ProcessTime:  processingTime,
		Error:        err,
	}
}

func collectResults(results <-chan ProcessResult, failedOrders chan<- Order, wg *sync.WaitGroup, totalOrders int) {
	defer wg.Done()

	var successful, failed int

	for i := 0; i < totalOrders; i++ {
		result := <-results

		if result.Success {
			fmt.Printf("Zamówienie #%d od %s przetworzone pomyślnie w czasie %v\n",
				result.OrderID, result.CustomerName, result.ProcessTime)
			successful++
		} else {
			fmt.Printf("Zamówienie #%d od %s NIE powiodło się: %v\n",
				result.OrderID, result.CustomerName, result.Error)
			failed++

			order := Order{
				ID:           result.OrderID,
				CustomerName: result.CustomerName,
			}
			failedOrders <- order
		}
	}

	close(failedOrders)

	fmt.Println("\n--- STATYSTYKI PRZETWARZANIA ---")
	fmt.Printf("Łącznie przetworzono: %d zamówień\n", totalOrders)
	fmt.Printf("Udane: %d (%.1f%%)\n", successful, float64(successful)/float64(totalOrders)*100)
	fmt.Printf("Nieudane: %d (%.1f%%)\n", failed, float64(failed)/float64(totalOrders)*100)
}

func retryFailedOrders(failedOrders <-chan Order, results chan<- ProcessResult, wg *sync.WaitGroup) {
	defer wg.Done()

	var retryCount int

	for order := range failedOrders {
		retryCount++
		fmt.Printf("Ponowna próba dla zamówienia #%d od %s\n", order.ID, order.CustomerName)

		var err error
		success := true

		if rand.IntN(20) == 0 { //5% szansy na błąd
			err = errors.New("błąd ponownego przetwarzania zamówienia")
			success = false
		}

		processingTime := time.Duration(rand.IntN(200)+100) * time.Millisecond
		time.Sleep(processingTime)

		result := ProcessResult{
			OrderID:      order.ID,
			CustomerName: order.CustomerName,
			Success:      success,
			ProcessTime:  processingTime,
			Error:        err,
		}

		results <- result
	}

	fmt.Printf("Zakończono ponowne próby dla %d zamówień\n", retryCount)
}
