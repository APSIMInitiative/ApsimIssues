package main

import (
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// Creates a scatter plot with the given parameters
func createLinePlot(dates []time.Time, y []int, title, xlabel, ylabel, fileName string, series ...map[time.Time]int) {
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

	line, err := plotter.NewLine(getXYPairs(dates, y))
	if err != nil {
		panic(err)
	}
	p.Add(line)

	if series != nil && len(series) > 0 {
		for i := 0; i < len(series); i++ {
			xn := sortKeys(series[i])
			var yn []int
			for _, x := range xn {
				yn = append(yn, series[i][x])
			}
			line, err = plotter.NewLine(getXYPairs(xn, yn))
			if err != nil {
				panic(err)
			}
			p.Add(line)
		}
	}

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
