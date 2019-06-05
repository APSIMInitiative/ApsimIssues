package main

import (
	"fmt"
	"time"

	"github.com/octokit/go-octokit/octokit"
)

// Increments all values in the map whose date key lies on or after
// a given date.
func incrementAfterDate(issues *map[time.Time]int, date time.Time) {
	for key := range *issues {
		if key.After(date) || key == date {
			(*issues)[key]++
		}
	}
}

// Decrements all values in the map whose date key lies on or after
// a given date.
func decrementAfterDate(issues *map[time.Time]int, date time.Time) {
	for key := range *issues {
		if key.After(date) || key == date {
			(*issues)[key]--
		}
	}
}

// Gets a map of dates to the number of open issues on that date.
func getOpenIssuesByDate(issues []octokit.Issue) map[time.Time]int {
	issuesByDate := make(map[time.Time]int)
	// Initialise the map with value for each date set to 0.
	for _, issue := range issues {
		issuesByDate[issue.CreatedAt] = 0
	}

	for _, issue := range issues {
		incrementAfterDate(&issuesByDate, issue.CreatedAt)
		if issue.ClosedAt != nil {
			decrementAfterDate(&issuesByDate, *issue.ClosedAt)
		}
	}
	return issuesByDate
}

// Gets a map of dates to the number of issues opened on or before that
// date.
func getCumOpenIssuesByDate(issues []octokit.Issue) map[time.Time]int {
	issuesByDate := make(map[time.Time]int)
	// Initialise the map with value for each date set to 0.
	for _, issue := range issues {
		issuesByDate[issue.CreatedAt] = 0
	}

	for _, issue := range issues {
		incrementAfterDate(&issuesByDate, issue.CreatedAt)
	}
	return issuesByDate
}

// Gets a map of dates to the number of closed issuse on that date.
func getCumIssuesClosedByDate(issues []octokit.Issue) map[time.Time]int {
	closed := make(map[time.Time]int)
	// Initialise the map with value for each date set to 0.
	for _, issue := range issues {
		closed[issue.CreatedAt] = 0
		if issue.ClosedAt != nil {
			closed[*issue.ClosedAt] = 0
		}
	}

	for _, issue := range issues {
		if issue.ClosedAt != nil {
			incrementAfterDate(&closed, *issue.ClosedAt)
		}
	}
	return closed
}

// Graphs the number of open bugs over time.
func graphIssuesByDate(issues []octokit.Issue, graphFileName string) {
	// Generate a map of issues over time.
	issuesOpenedByDate := getOpenIssuesByDate(issues)

	title := "Change in number of open bugs over time"

	createLinePlot(
		title,
		"Date",
		"Number of open bugs",
		graphFileName,
		seriesFromMap(title, issuesOpenedByDate))
	fmt.Printf("Generated graph '%s'\n", graphFileName)
}

func graphOpenedVsClosed(issues []octokit.Issue, graphFileName string) {
	// Generate a map of issues over time.
	opened := seriesFromMap("Total issues opened",
		getCumOpenIssuesByDate(issues))
	closed := seriesFromMap("Total issues closed",
		getCumIssuesClosedByDate(issues))

	createLinePlot(
		"Total issues opened and closed over time",
		"Date",
		"Number of open bugs",
		graphFileName,
		opened,
		closed)
	fmt.Printf("Generated graph '%s'\n", graphFileName)
}
