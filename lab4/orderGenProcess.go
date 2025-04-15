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
	var failedOrderIDs = make(map[int]bool)
	
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
			failedOrderIDs[result.OrderID] = true

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
	var retrySuccessful, retryFailed int
	var retryResultsMutex sync.Mutex
	var retriedOrders []Order

	for order := range failedOrders {
		retriedOrders = append(retriedOrders, order)
	}
	
	retryCount = len(retriedOrders)
	
	if retryCount > 0 {
		var retryWg sync.WaitGroup
		
		for _, order := range retriedOrders {
			retryWg.Add(1)
			go func(o Order) {
				defer retryWg.Done()
				
				fmt.Printf("Ponowna próba dla zamówienia #%d od %s\n", o.ID, o.CustomerName)

				var err error
				success := true

				if rand.IntN(20) == 0 { // 5% szansy na błąd
					err = errors.New("błąd ponownego przetwarzania zamówienia")
					success = false
				}

				processingTime := time.Duration(rand.IntN(200)+100) * time.Millisecond
				time.Sleep(processingTime)

				result := ProcessResult{
					OrderID:      o.ID,
					CustomerName: o.CustomerName,
					Success:      success,
					ProcessTime:  processingTime,
					Error:        err,
				}

				retryResultsMutex.Lock()
				if success {
					retrySuccessful++
				} else {
					retryFailed++
				}
				retryResultsMutex.Unlock()
				
				results <- result
			}(order)
		}
		
		retryWg.Wait()
		
		fmt.Printf("Zakończono ponowne próby dla %d zamówień\n", retryCount)
		
		fmt.Println("\n--- STATYSTYKI PONOWNYCH PRÓB ---")
		fmt.Printf("Łącznie ponownych prób: %d\n", retryCount)
		fmt.Printf("Udane ponowne próby: %d (%.1f%%)\n", retrySuccessful, 
			float64(retrySuccessful)/float64(retryCount)*100)
		fmt.Printf("Nieudane ponowne próby: %d (%.1f%%)\n", retryFailed, 
			float64(retryFailed)/float64(retryCount)*100)
	}
}