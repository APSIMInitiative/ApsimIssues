package main

import (
	"fmt"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Creates a scatter plot with the given parameters
func createLinePlot(title, xlabel, ylabel, fileName string, data ...series) {
	if data == nil || len(data) < 1 {
		panic("Graph error: No series provided")
	}
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	var baseFontSize vg.Length = 36

	// Title formatting
	p.Title.Text = title
	p.Title.Font.Size = 48

	// x-axis formatting
	p.X.Label.Text = xlabel
	p.X.Label.Font.Size = baseFontSize
	if _, ok := data[0].(dateSeries); ok {
		p.X.Tick.Marker = plot.TimeTicks{
			Format: "Jan 2006",
			Ticker: TimeTicker{},
		}
	}
	p.X.Tick.Label.Font.Size = baseFontSize
	p.Y.Tick.Label.Font.Size = baseFontSize

	// y-axis formatting
	p.Y.Label.Text = ylabel
	p.Y.Label.Font.Size = 42

	// Legend formatting
	p.Legend.Font.Size = baseFontSize
	p.Legend.ThumbnailWidth = 24
	p.Legend.Top = true
	p.Legend.Left = true

	colours := palette.Rainbow(len(data)+1, 0, 1, 1, 1, 1).Colors()

	for i := 0; i < len(data); i++ {
		line, err := plotter.NewLine(data[i].getXYPairs())
		line.LineStyle.Width = 2
		if len(data) > 1 {
			line.LineStyle.Color = colours[i]
			p.Legend.Add(data[i].getName(), line)
		}
		if err != nil {
			panic(err)
		}
		p.Add(line)
	}

	// Write to disk
	err = p.Save(1920, 1080, fileName)
	if err != nil {
		panic(err)
	}
	if !settings.Quiet {
		fmt.Printf("Generated graph '%s'\n", fileName)
	}
}

func createBarChart(title string, xAxisLabel string, yAxisLabel string, fileName string, data ...barSeries) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	var baseFontSize vg.Length = 36
	var barLength vg.Length = 18

	// Title formatting
	p.Title.Text = title
	p.Title.Font.Size = 48

	// x-axis formatting
	p.X.Label.Text = xAxisLabel
	p.X.Label.Font.Size = baseFontSize
	p.X.Tick.Label.Font.Size = baseFontSize
	p.Y.Tick.Label.Font.Size = baseFontSize

	// y-axis formatting
	p.Y.Label.Text = yAxisLabel
	p.Y.Label.Font.Size = 42

	// Legend formatting
	p.Legend.Font.Size = baseFontSize
	p.Legend.ThumbnailWidth = 24
	p.Legend.Top = true
	p.Legend.Left = true

	colours := palette.Rainbow(len(data)+1, 0, 1, 1, 1, 1).Colors()

	var prevChart *plotter.BarChart
	for i, series := range data {
		chart, err := plotter.NewBarChart(series, barLength)
		if err != nil {
			panic(err)
		}
		p.Add(chart)
		chart.Color = colours[i]
		p.Legend.Add(series.Name, chart)
		if i > 0 {
			chart.StackOn(prevChart)
		}
		prevChart = chart
	}

	p.NominalX(data[0].Names...)

	// for i := 0; i < len(data); i++ {
	// 	line, err := plotter.NewLine(data[i].getXYPairs())
	// 	line.LineStyle.Width = 2
	// 	if len(data) > 1 {
	// 		line.LineStyle.Color = colours[i]
	// 		p.Legend.Add(data[i].getName(), line)
	// 	}
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	p.Add(line)
	// }

	// Write to disk
	err = p.Save(2560, 1440, fileName)
	if err != nil {
		panic(err)
	}
	if !settings.Quiet {
		fmt.Printf("Generated graph '%s'\n", fileName)
	}
}

// Ticks generates ticks for the time axis
func (T TimeTicker) Ticks(min, max float64) (ticks []plot.Tick) {
	minTime := time.Unix(int64(min), 0)
	maxTime := time.Unix(int64(max), 0)

	years, months, _, _, _, _ := diff(maxTime, minTime)
	totalMonths := months + years*12
	// We want to show up to a maximum of 10 ticks on the x axis.
	var increment int
	maxNoTicks := 10
	if totalMonths <= maxNoTicks {
		increment = 1
	} else if totalMonths/3 <= maxNoTicks {
		increment = 3
	} else if totalMonths/6 <= maxNoTicks {
		increment = 6
	} else if totalMonths/12 <= maxNoTicks {
		increment = 12
	} else {
		// If we are plotting more than 20 years of data, just use the
		// default gonum time tick algorithm.
		return plot.TimeTicks{}.Ticks(min, max)
	}
	for i := minTime; i.Before(maxTime) || i == maxTime; i = i.AddDate(0, increment, 0) {
		tick := plot.Tick{
			Value: float64(i.Unix()),
			Label: i.Format("Jan 2006"),
		}
		ticks = append(ticks, tick)
	}
	return
}

// TimeTicker implements plot.Ticker.
type TimeTicker struct {
}
