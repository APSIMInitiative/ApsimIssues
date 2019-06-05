package main

import (
	"fmt"
	"time"

	"gonum.org/v1/plot/plotter"
)

type series struct {
	X    []time.Time
	Y    []int
	Name string
}

func seriesFromMap(title string, data map[time.Time]int) series {
	var result series
	result.Name = title
	result.X = sortKeys(data)
	for _, date := range result.X {
		result.Y = append(result.Y, data[date])
	}
	return result
}

func (S series) getXYPairs() plotter.XYs {
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
