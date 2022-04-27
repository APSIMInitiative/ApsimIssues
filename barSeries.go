package main

import (
	"github.com/octokit/go-octokit/octokit"
	"gonum.org/v1/plot/plotter"
)

// barSeries represents a set of bars to go on a bar chart.
type barSeries struct {
	Name   string
	Names  []string
	Values []float64
}

// Len returns the number of values.
func (S barSeries) Len() int {
	return len(S.Names)
}

// Value returns a value.
func (S barSeries) Value(index int) float64 {
	return S.Values[index]
}

func barSeriesFromGroups(name string, groups map[string][]octokit.Issue) barSeries {
	series := barSeries{}
	series.Name = name

	for user, issues := range groups {
		series.Names = append(series.Names, user)
		series.Values = append(series.Values, float64(len(issues)))
	}

	return series
}

func getXY(s barSeries) plotter.XY {
	return plotter.XY{X: s.Values[0], Y: s.Values[0]}
}
