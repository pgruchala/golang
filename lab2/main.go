package main

import (
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"sort"
)

type Entry struct {
	Code string
	DELabel   string
    ENLabel   string
    ESLabel   string
    FRLabel   string
    PTLabel   string
    ShortCode string
}

func main() {
	file,err := os.Open("nomenclature-cpv.csv")
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'
	records,err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var items []Entry
	for _,record := range records[1:] {
		item := Entry{
			Code: record[0],
			DELabel: record[1],
			ENLabel: record[2],
			ESLabel: record[3],
			FRLabel: record[4],
			PTLabel: record[5],
			ShortCode: record[6],
		}
		items = append(items, item)
	}
	fmt.Println("Sort by Code")
	sort.Slice(items, func(i, j int) bool {
		return items[i].Code < items[j].Code
	})
	for _,item := range items {
		fmt.Println(item)
	}
	fmt.Println("-------------------------------------------------")
	fmt.Println("Sort by ENLabel")
	sort.Slice(items, func(i, j int) bool {
		return items[i].ENLabel < items[j].ENLabel
	})
	for _,item := range items {
		fmt.Println(item)
	}

	fmt.Println("-------------------------------------------------")
	fmt.Println("Średnia długość nazwy dla każdego języka")
	var sumDE,sumEN,sumES,sumFR,sumPT int
	for _,item := range items {
		sumDE += len(item.DELabel)
		sumEN += len(item.ENLabel)
		sumES += len(item.ESLabel)
		sumFR += len(item.FRLabel)
		sumPT += len(item.PTLabel)
	}
	fmt.Println("DE: ",float64(sumDE)/float64(len(items)))
	fmt.Println("EN: ",float64(sumEN)/float64(len(items)))
	fmt.Println("ES: ",float64(sumES)/float64(len(items)))
	fmt.Println("FR: ",float64(sumFR)/float64(len(items)))
	fmt.Println("PT: ",float64(sumPT)/float64(len(items)))

}