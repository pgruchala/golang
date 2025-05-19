package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type CSVLoader struct{}

func NewCSVLoader() *CSVLoader {
	return &CSVLoader{}
}

func (l *CSVLoader) LoadData(filePath string) ([]Entry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("błąd otwarcia pliku %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("błąd przy odczycie pliku CSV: %w", err)
	}
	var entries []Entry
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Błąd odczytu linii %v: %v", record, err)
			continue
		}
		if len(record) < 6 {
			fmt.Printf("Błąd odczytu linii %v: Niekompletna linia", record)
			continue
		}
		dateFormated := record[0]
		closingFormated := strings.ReplaceAll(record[1], "$", "")
		volumeFormatted := strings.ReplaceAll(record[2], ",", "")
		openingFormated := strings.ReplaceAll(record[3], "$", "")
		highFormated := strings.ReplaceAll(record[4], "$", "")
		lowFormated := strings.ReplaceAll(record[5], "$", "")
		date, err := time.Parse("01/02/2006", dateFormated)
		if err != nil {
			fmt.Printf("Błąd parsowania daty '%s': %v. Pomijanie linii.\n", dateFormated, err)
			continue
		}

		closePrice, err := strconv.ParseFloat(closingFormated, 64)
		if err != nil {
			fmt.Printf("Błąd parsowania ceny zamknięcia '%s': %v. Pomijanie linii.\n", closingFormated, err)
			continue
		}

		volume, err := strconv.ParseUint(volumeFormatted, 10, 64)
		if err != nil {
			fmt.Printf("Błąd parsowania wolumenu '%s': %v. Pomijanie linii.\n", volumeFormatted, err)
			continue
		}

		openPrice, err := strconv.ParseFloat(openingFormated, 64)
		if err != nil {
			fmt.Printf("Błąd parsowania ceny otwarcia '%s': %v. Pomijanie linii.\n", openingFormated, err)
			continue
		}

		highPrice, err := strconv.ParseFloat(highFormated, 64)
		if err != nil {
			fmt.Printf("Błąd parsowania ceny najyższej '%s': %v. Pomijanie linii.\n", highFormated, err)
			continue
		}

		lowPrice, err := strconv.ParseFloat(lowFormated, 64)
		if err != nil {
			fmt.Printf("Błąd parsowania ceny najniższej '%s': %v. Pomijanie linii.\n", lowFormated, err)
			continue
		}
		entries = append(entries, Entry{
			Date:   date,
			Close:  closePrice,
			Volume: float64(volume),
			Open:   openPrice,
			High:   highPrice,
			Low:    lowPrice,
		})
	}
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	return entries, nil

}
