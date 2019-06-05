package main

import (
	"image/color"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Creates a scatter plot with the given parameters
func createScatterPlot(dates []time.Time, y []int, title, xlabel, ylabel, fileName string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = title
	p.Title.Font.Size = 32
	p.X.Label.Text = xlabel
	p.X.Tick.Marker = plot.TimeTicks{Format: "Jan 2006"}
	p.Y.Label.Text = ylabel

	scatter, err := plotter.NewScatter(getXYPairs(dates, y))
	if err != nil {
		panic(err)
	}

	scatter.GlyphStyle.Color = color.Black
	p.Add(scatter)
	// Write to disk
	err = p.Save(1920, 1080, fileName)
	if err != nil {
		panic(err)
	}
}

func getXYPairs(dates []time.Time, y []int) plotter.XYs {
	if len(dates) != len(y) {
		panic("Error: x/y data length mismatch")
	}
	points := make(plotter.XYs, len(dates))
	for i, date := range dates {
		points[i].X = float64(date.Unix())
		points[i].Y = float64(y[i])
	}
	return points
}
