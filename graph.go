package main

import (
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
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
	p.Title.Text = title
	p.Title.Font.Size = 32
	p.X.Label.Text = xlabel
	p.X.Label.Font.Size = 20
	p.X.Tick.Marker = plot.TimeTicks{
		Format: "Jan 2006",
		Ticker: TimeTicker{},
	}

	p.Y.Label.Text = ylabel
	p.Y.Label.Font.Size = 20
	p.Legend.Font.Size = 24
	p.Legend.ThumbnailWidth = 24
	colours := palette.Rainbow(len(data)+1, 0, 1, 1, 1, 1).Colors()

	for i := 0; i < len(data); i++ {
		line, err := plotter.NewLine(data[i].getXYPairs())
		if len(data) > 1 {
			line.LineStyle.Color = colours[i]
			p.Legend.Add(data[i].Name, line)
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
}

// Ticks generates ticks for the time axis
func (T TimeTicker) Ticks(min, max float64) (ticks []plot.Tick) {
	minTime := time.Unix(int64(min), 0)
	maxTime := time.Unix(int64(max), 0)

	years, months, _, _, _, _ := diff(maxTime, minTime)
	totalMonths := months + years*12
	// We want to show up to a maximum of 20 ticks on the x axis.
	var increment int
	if totalMonths <= 20 {
		increment = 1
	} else if totalMonths/3 <= 20 {
		increment = 3
	} else if totalMonths/6 <= 20 {
		increment = 6
	} else if totalMonths/12 <= 20 {
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
