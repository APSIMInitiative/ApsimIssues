package main

import (
	"fmt"
	"time"

	"gonum.org/v1/plot/plotter"
)

// dateSeries is an implementation of series which can be used to graph
// time.Time data on the x-axis against ints on the y-axis.
type dateSeries struct {
	X    []time.Time
	Y    []int
	Name string
}

// seriesFromMap generates a series from a map of times to ints.
func seriesFromMap(title string, data map[time.Time]int) dateSeries {
	var result dateSeries
	result.Name = title
	result.X = sortKeys(data)
	for _, date := range result.X {
		result.Y = append(result.Y, data[date])
	}
	return result
}

// getXYPairs returns the graph's XY data as a plotter.XYs object.
func (S dateSeries) getXYPairs() plotter.XYs {
	if len(S.X) != len(S.Y) {
		panic(fmt.Sprintf("Error in series '%s': x/y data length mismatch", S.Name))
	}
	points := make(plotter.XYs, len(S.X))
	for i, date := range S.X {
		points[i].X = float64(date.Unix())
		points[i].Y = float64(S.Y[i])
	}
	return points
}

// getName returns the series' name to be used in the legend.
func (S dateSeries) getName() string {
	return S.Name
}
