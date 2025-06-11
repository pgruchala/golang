package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func loadData() []City {
	file, err := os.Open("worldcities.csv")
	if err != nil {
		fmt.Println("Error opening file", err)
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error reading from file", err)
		return nil
	}
	var cities []City
	for {
		record, err := reader.Read()
		if err == io.EOF {
			fmt.Println("Reached end of file")
			break
		}
		if err != nil {
			fmt.Println("Error loading city: ", err)
		}
		latitude, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Println("Invalid latitude", err)
			continue
		}
		longitude, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Println("Invalid longitude", err)
			continue
		}
		city := City{
			Name:      record[1],
			Latitude:  latitude,
			Longitude: longitude,
		}
		cities = append(cities, city)
	}
	return cities
}
