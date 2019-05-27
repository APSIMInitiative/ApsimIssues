package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"image/color"
	"math/rand"
)

func createGraph(pulls []pullRequest, fileName string) {
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
			pts[i].X = pts[i - 1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10 * rand.Float64()
	}
	return pts
}