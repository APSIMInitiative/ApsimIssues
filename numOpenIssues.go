package main

import (
	"time"

	"github.com/octokit/go-octokit/octokit"
	//"time"
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

func numIssuesByDate(client *octokit.Client, owner, repo string) {
	var issues []octokit.Issue
	if fileExists(cache) {
		issues = issuesFromCache(cache)
	} else {
		// Get all issues (open + closed) on ApsimX.
		issues = getAllIssues(client, owner, repo, !quiet)
	}

	// Generate a map of issues over time.
	issuesOpenedByDate := getIssuesOpenedByDate(issues)

	// Convert the map into two arrays.
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
		"openIssues.png")

	// Cache results for next time.
	writeToCache(cache, issues)
}
