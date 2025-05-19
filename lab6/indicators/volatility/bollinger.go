package volatility

import (
	"fmt"
	"lab6/indicators/trend"
	"math"
)

type BollingerBandPoint struct {
	MiddleBand float64
	UpperBand  float64
	LowerBand  float64
}

func CalculateBollinger(prices []float64, period int, stdDev float64) ([]BollingerBandPoint, error) {
	if period <= 0 {
		return nil, fmt.Errorf("okres musi być dodatni")
	}
	if stdDev <= 0 {
		return nil, fmt.Errorf("liczba odchyleń standardowych musi być dodatnia")
	}
	if len(prices) < period {
		results := make([]BollingerBandPoint, len(prices))
		for i := range results {
			results[i] = BollingerBandPoint{math.NaN(), math.NaN(), math.NaN()}
		}
		return results, nil
	}

	smaValues, err := trend.CalculateSMA(prices, period)
	if err != nil {
		return nil, fmt.Errorf("błąd obliczania SMA dla Wstęg Bollingera: %w", err)
	}

	bbPoints := make([]BollingerBandPoint, len(prices))

	for i := 0; i < period-1; i++ {
		bbPoints[i] = BollingerBandPoint{math.NaN(), math.NaN(), math.NaN()}
	}

	smaOffset := len(prices) - len(smaValues)

	for i := period - 1; i < len(prices); i++ {

		smaIndex := i - smaOffset
		
		if smaIndex >= len(smaValues) || smaIndex < 0 {
			bbPoints[i] = BollingerBandPoint{math.NaN(), math.NaN(), math.NaN()}
			continue
		}
		
		if math.IsNaN(smaValues[smaIndex]) {
			bbPoints[i] = BollingerBandPoint{math.NaN(), math.NaN(), math.NaN()}
			continue
		}

		var sumSqDiff float64
		startIdx := i - period + 1
		if startIdx < 0 {
			bbPoints[i] = BollingerBandPoint{math.NaN(), math.NaN(), math.NaN()}
			continue
		}

		for j := startIdx; j <= i; j++ {
			diff := prices[j] - smaValues[smaIndex]
			sumSqDiff += diff * diff
		}
		stdDev2 := math.Sqrt(sumSqDiff / float64(period))

		bbPoints[i].MiddleBand = smaValues[smaIndex]
		bbPoints[i].UpperBand = smaValues[smaIndex] + (stdDev * stdDev2)
		bbPoints[i].LowerBand = smaValues[smaIndex] - (stdDev * stdDev2)
	}

	return bbPoints, nil
}
