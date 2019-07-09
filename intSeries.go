package main

import (
	"fmt"

	"gonum.org/v1/plot/plotter"
)

// intSeries is an implementation of the series interface which plots
// ints against ints.
type intSeries struct {
	X    []int
	Y    []int
	Name string
}

// getXYPairs returns the graph's XY data as a plotter.XYs object.
func (s intSeries) getXYPairs() plotter.XYs {
	if len(s.X) != len(s.Y) {
		panic(fmt.Sprintf("Error in series '%s': x/y data length mismatch", s.Name))
	}
	points := make(plotter.XYs, len(s.X))
	for i, x := range s.X {
		points[i].X = float64(x)
		points[i].Y = float64(s.Y[i])
	}
	return points
}

// getName returns the series' name to be used in the legend.
func (s intSeries) getName() string {
	return s.Name
}

func createIntSeries(x, y dateSeries, name string) intSeries {
	if len(x.Y) == len(y.Y) {
		return intSeries{
			X:    x.Y,
			Y:    y.Y,
			Name: name,
		}
	}
	var xData []int
	var yData []int
	for i, date := range x.X {
		j := indexOf(y.X, date)
		if j >= 0 {
			xData = append(xData, x.Y[i])
			yData = append(yData, y.Y[j])
		}
	}
	return intSeries{
		X:    xData,
		Y:    yData,
		Name: name,
	}
}
