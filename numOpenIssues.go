package main

import (
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

func getIssuesOpenedByDate(issues []octokit.Issue) map[time.Time]int {
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

func graphIssuesByDate(issues []octokit.Issue, graphFileName string) {
	// Generate a map of issues over time.
	issuesOpenedByDate := getIssuesOpenedByDate(issues)

	// Convert the map into two arrays. Maps have no concept of order,
	// but the arrays need be ordered by date ascending.
	dates := sortKeys(issuesOpenedByDate)
	var numIssues []int
	for _, date := range dates {
		numIssues = append(numIssues, issuesOpenedByDate[date])
	}

	createScatterPlot(
		dates,
		numIssues, "Change in number of open bugs over time",
		"Date",
		"Number of open bugs",
		graphFileName)
}
