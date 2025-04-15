package main

import (
	"fmt"
	"sync"
)

func main() {
	numOrders := 20

	numWorkers := 5
	
	ordersChan := make(chan Order, numOrders)
	resultsChan := make(chan ProcessResult, numOrders*2) // 2 razy bo mogą być ponowne próby
	failedOrdersChan := make(chan Order, numOrders)
	
	var wg sync.WaitGroup
	
	go orderGenerator(numOrders, ordersChan)
	
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, ordersChan, resultsChan, &wg)
	}
	
	wg.Add(1)
	go collectResults(resultsChan, failedOrdersChan, &wg, numOrders)
	
	wg.Add(1)
	go retryFailedOrders(failedOrdersChan, resultsChan, &wg)
	
	wg.Wait()
	
	fmt.Println("\nSymulacja zakończona")
}