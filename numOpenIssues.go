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
		//closed[issue.CreatedAt] = 0
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

func graphOpenedVsClosedForUser(issues []octokit.Issue, pulls []octokit.PullRequest, userName, graphFileName string) {
	// Get data for issues fixed by the user.
	dates, fixed := getBugFixRate(pulls, userName)
	fixedSeries := series{
		Name: fmt.Sprintf("Total fixed by %s", userName),
		X:    dates,
		Y:    fixed,
	}

	if len(dates) < 1 {
		return
	}

	// We only want to graph data on or after the date of the first bug fixed by the user.
	dateFirstBugfix := dates[0]
	openedAfterDate := filterIssues(issues, func(issue octokit.Issue) bool {
		return issue.CreatedAt.After(dateFirstBugfix) || issue.CreatedAt == dateFirstBugfix
	})
	closedAfterDate := filterIssues(issues, func(issue octokit.Issue) bool {
		return issue.ClosedAt != nil &&
			((*issue.ClosedAt).After(dateFirstBugfix) || *issue.ClosedAt == dateFirstBugfix)
	})

	// Generate a map of cumulative issues opened and closed over time.
	opened := seriesFromMap("Total issues opened", getCumOpenIssuesByDate(openedAfterDate))
	closed := seriesFromMap("Total issues closed", getCumIssuesClosedByDate(closedAfterDate))

	createLinePlot(
		fmt.Sprintf("Total issues opened and closed over time since %s's first bugfix", userName),
		"Date",
		"Number of open bugs",
		graphFileName,
		opened,
		closed,
		fixedSeries)
	fmt.Printf("Generated graph '%s'\n", graphFileName)
}

func graphOpenedVsClosedForUsers(issues []octokit.Issue, pulls []octokit.PullRequest, graphFileName string, users ...string) {
	// Get data for issues fixed for each user.
	var userSeries []series
	for _, userName := range users {
		dates, fixed := getBugFixRate(pulls, userName)
		userSeries = append(userSeries, series{
			Name: fmt.Sprintf("Total fixed by %s", userName),
			X:    dates,
			Y:    fixed,
		})
	}
	// Generate a map of cumulative issues opened and closed over time.
	opened := seriesFromMap("Total issues opened",
		getCumOpenIssuesByDate(issues))
	allSeries := append(userSeries, opened)

	closed := seriesFromMap("Total issues closed",
		getCumIssuesClosedByDate(issues))
	allSeries = append(allSeries, closed)

	createLinePlot(
		"Total issues opened and closed over time",
		"Date",
		"Number of bugs",
		graphFileName,
		allSeries...)
	fmt.Printf("Generated graph '%s'\n", graphFileName)
}

// graphBugfixRateByUser graphs the cumulative issues opened and closed
// vs the total number of issues fixed by each user who has fixed at
// least `minN` bugs.
func graphBugfixRateByUser(issues []octokit.Issue, pulls []octokit.PullRequest, graphFileName string, minN int) {
	// Get data for issues fixed for each user.
	var userSeries []series
	dataByUser := pullsGroupedByUser(pulls)
	for user := range dataByUser {
		// Generate a map of dates to number of issues referenced in pull requests.
		issuesByDate := getCumIssuesByDate(dataByUser[user])

		// Only graph data for this user if they have fixed at least `minN` bugs.
		numBugs := issuesByDate[getLastDate(issuesByDate)]
		if numBugs >= minN {
			seriesTitle := user
			userSeries = append(userSeries, seriesFromMap(seriesTitle, issuesByDate))
		}
	}

	// Generate a map of cumulative issues opened and closed over time.
	opened := seriesFromMap("Total issues opened",
		getCumOpenIssuesByDate(issues))
	userSeries = append(userSeries, opened)

	closed := seriesFromMap("Total issues closed",
		getCumIssuesClosedByDate(issues))
	userSeries = append(userSeries, closed)

	createLinePlot(
		fmt.Sprintf("Bugs fixed over time for all users who have fixed at least %d bugs", minN),
		"Date",
		"Number of bugs",
		graphFileName,
		userSeries...)
	fmt.Printf("Generated graph '%s'\n", graphFileName)
}
