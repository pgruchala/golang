package main

import (
	"math/rand/v2"
)

type Item struct {
	name string
	price float64
}

func generateNewOrder(id int) Order {
	produkty := []Item{{"chleb",5.50},{"bułki",2.50},{"dżemy",8.00},{"jabłka",4.20},{"LEGO",69.99},{"sos totkowy",9.20}, {"ryba",18.40}, {"kebab",25.00}, {"dziecko",25000.00},{"samolot",40.42},}
	customerName := []string{"James","Michael","Robert","John","David","WIlliam","RIchard","Joseph","Thomas","Christopher","Mary","Patricia","Jennifer","Linda","Elizabeth","Barbara","Susan","Jessica","Karen","Sarah"}
	amount := rand.IntN(15)
	var items = []string{}
	var total float64
	for i := 0; i < amount; i++ {
		randomProductIndex := rand.IntN(10)
		selectedItem := produkty[randomProductIndex]

		items = append(items, selectedItem.name)
		total += selectedItem.price 
	}
	randomName := customerName[rand.IntN(len(customerName))]
	return Order{ID: id,CustomerName: randomName,Items: items, TotalAmount: total}
}

