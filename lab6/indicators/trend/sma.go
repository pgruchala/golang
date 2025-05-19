package trend

import (
	"fmt"
	_ "lab6/data"
	"math"
)

func CalculateSMA(prices []float64, period int) ([]float64, error) {
	if period < 0 {
		return nil, fmt.Errorf("okres nie może być ujemny")
	}
	if len(prices) < period {
		result := make([]float64, len(prices))
		for i := range result {
			result[i] = math.NaN()
		}
		return result, nil
	}
	smaValues := make([]float64, len(prices)-period+1)
	for i := 0; i <= len(prices)-period; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i+j]
		}
		smaValues[i] = sum / float64(period)
	}
	return smaValues, nil

}
