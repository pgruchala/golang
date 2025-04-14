package main

import "time"

// Struktura modelująca zamówienie
type Order struct {
	ID           int
	CustomerName string
	Items        []string
	TotalAmount  float64
}

// Struktura modelująca prze
type ProcessResult struct {
	OrderID      int
	CustomerName string
	Success      bool
	ProcessTime  time.Duration
	Error        error
}
