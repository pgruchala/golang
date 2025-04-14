package main

import (
	"fmt"
	"sync"
)

func worker(id int, orders <- chan Order, results chan <- ProcessResult, wg *sync.WaitGroup){
	defer wg.Done()

	for order := range orders {
		fmt.Printf("Pracownik ID%d przetwarza zamównienie o id %d od %s\n",id, order.ID, order.CustomerName)
		result := processOrder(order)
		results <- result
	}
}