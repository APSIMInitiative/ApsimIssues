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

func graphIssuesByDate(issues []octokit.Issue, graphFileName string) {
	// Generate a map of issues over time.
	issuesOpenedByDate := getOpenIssuesByDate(issues)

	// Convert the map into two arrays. Maps have no concept of order,
	// but the arrays need be ordered by date ascending.
	dates := sortKeys(issuesOpenedByDate)
	var numIssues []int
	for _, date := range dates {
		numIssues = append(numIssues, issuesOpenedByDate[date])
	}

	createLinePlot(
		dates,
		numIssues,
		"Change in number of open bugs over time",
		"Date",
		"Number of open bugs",
		graphFileName)
	fmt.Printf("Generated graph '%s'\n", graphFileName)
}

func graphOpenedVsClosed(issues []octokit.Issue, graphFileName string) {
	// Generate a map of issues over time.
	issuesOpenedByDate := getCumOpenIssuesByDate(issues)

	// Convert the map into two arrays. Maps have no concept of order,
	// but the arrays need be ordered by date ascending.
	dates := sortKeys(issuesOpenedByDate)
	var numIssues []int
	for _, date := range dates {
		numIssues = append(numIssues, issuesOpenedByDate[date])
	}

	closedByDate := getCumIssuesClosedByDate(issues)

	createLinePlot(
		dates,
		numIssues,
		"Change in number of open bugs over time",
		"Date",
		"Number of open bugs",
		graphFileName,
		closedByDate)
}
