package data

import "time"

type Entry struct {
	Date   time.Time
	Close  float64
	Volume float64
	Open   float64
	High   float64
	Low    float64
}

type DataLoader interface {
	LoadData(filePath string) ([]Entry, error)
}
