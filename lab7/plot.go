package main

import (
	"image/color"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func generatePlot(data WeatherResponse, filename string) error {
	p := plot.New()

	p.Title.Text = "Wykres temperatur"
	p.X.Label.Text = "Data"
	p.Y.Label.Text = "Temperatura (°C)"

	maxTemps := make(plotter.XYs, len(data.Daily.Time))
	minTemps := make(plotter.XYs, len(data.Daily.Time))

	dates := make([]string, len(data.Daily.Time))
	for i := range data.Daily.Time {
		maxTemps[i].X = float64(i)
		maxTemps[i].Y = data.Daily.TempMax[i]
		minTemps[i].X = float64(i)
		minTemps[i].Y = data.Daily.TempMin[i]
		t, _ := time.Parse("2006-01-02", data.Daily.Time[i])
		dates[i] = t.Format("02-01") // Format DD-MM
	}

	p.NominalX(dates...)

	maxLine, err := plotter.NewLine(maxTemps)
	if err != nil {
		return err
	}
	maxLine.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	maxLine.LineStyle.Width = vg.Points(2)

	minLine, err := plotter.NewLine(minTemps)
	if err != nil {
		return err
	}
	minLine.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	minLine.LineStyle.Width = vg.Points(2)

	p.Add(maxLine, minLine)
	p.Legend.Add("Temp. Max", maxLine)
	p.Legend.Add("Temp. Min", minLine)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		return err
	}
	return nil
}
