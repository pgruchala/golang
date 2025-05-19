package momentum

import (
	"fmt"
	"math"
)



func CalculateCCI(highs []float64, lows []float64, closes []float64, period int) ([]float64, error) {
	if period <= 0 {
		return nil, fmt.Errorf("okres nie może być ujemny")
	}
	if len(highs) < period || len(lows) < period || len(closes) < period {
		return nil, fmt.Errorf("niewystarczająca ilość danych do obliczenia CCI dla okresu %d", period)
	}
	if len(highs) != len(lows) || len(highs) != len(closes) {
		return nil, fmt.Errorf("długości High, Low i Close muszą być takie same")
	}
	cciValues := make([]float64, len(highs)-period+1)
	for i := 0; i <= len(highs)-period; i++ {
		typicalPrices := make([]float64, period)
		for j := 0; j < period; j++ {
			typicalPrices[j] = (highs[i+j] + lows[i+j] + closes[i+j]) / 3
		}

		sumTP := 0.0
		for _, tp := range typicalPrices {
			sumTP += tp
		}
		meanTP := sumTP / float64(period)

		meanDeviation := 0.0
		for _, tp := range typicalPrices {
			meanDeviation += math.Abs(tp - meanTP)
		}
		meanDeviation = meanDeviation / float64(period)

		if meanDeviation == 0 {
			cciValues[i] = 0
		} else {
			cciValues[i] = (typicalPrices[period-1] - meanTP) / (0.015 * meanDeviation)
		}
	}
	return cciValues, nil

}
