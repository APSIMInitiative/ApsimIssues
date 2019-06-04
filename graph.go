package main

import (
	"image/color"
	"math/rand"
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

func createBugfixGraph(pulls []pullRequest, fileName string) {
	rand.Seed(int64(0))
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Cumulative bugs fixed over time"
	p.X.Label.Text = "Date"
	p.X.Tick.Marker = plot.TimeTicks{Format: "Jan 2006"}
	p.Y.Label.Text = "Total number of issues resolved"

	//err = plotutil.AddScatters(p, "", getXYs(pulls))
	scatter, err := plotter.NewScatter(getXYs(pulls))
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

func getXYs(pulls []pullRequest) plotter.XYs {
	points := make(plotter.XYs, len(pulls))
	issuesByDate := getIssuesByDate(pulls)
	sortedDates := sortKeys(issuesByDate)
	sum := 0
	for i, date := range sortedDates {
		points[i].X = float64(date.Unix())

		sum += issuesByDate[date]
		points[i].Y = float64(sum)
	}
	return points
}

func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}
