package main

import "gonum.org/v1/plot/plotter"

// series encapsulates a dataset which can be graphed.
type series interface {
	// getXYPairs returns the graph's XY data as a plotter.XYs object.
	getXYPairs() plotter.XYs

	// getName returns the series' name to be used in the legend.
	getName() string
}
