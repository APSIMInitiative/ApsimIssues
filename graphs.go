package main

import (
	"fmt"

	"github.com/octokit/go-octokit/octokit"
)

// graphBugFixRate graphs the cumulative number of bugs fixed by a user
// over time.
func graphBugFixRate(allPulls []octokit.PullRequest, username, graphFileName string) {
	bugFixRate := getBugFixRate(allPulls, username)
	title := fmt.Sprintf("Cumulative bugs fixed over time by %s", username)

	data := seriesFromMap(title, bugFixRate)

	createLinePlot(
		title,
		"Date",
		"Total Number of Issues Resolved",
		graphFileName,
		data)
	if bugFixRate != nil {
		fmt.Printf("%s has resolved %d issues.\n", username, bugFixRate[getLastDate(bugFixRate)])
	}
}

// graphIssuesByDate graphs the number of open bugs over time.
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
}

// graphOpenedVsClosed graphs two series:
// 1. Cumulative number of issues opened over time.
// 2. Cumulative number of issues closed over time.
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
}

// graphOpenedVsClosed graphs three series:
// 1. Cumulative number of issues opened over time.
// 2. Cumulative number of issues closed over time.
// 3. Cumulative number of issues fixed over time by a given user.
func graphOpenedVsClosedForUser(issues []octokit.Issue, pulls []octokit.PullRequest, userName, graphFileName string) {
	bugFixRate := getBugFixRate(pulls, userName)
	fixedSeries := seriesFromMap(
		fmt.Sprintf("Total fixed by %s", userName),
		bugFixRate)

	// We only want to graph data on or after the date of the first bug fixed by the user.
	dateFirstBugfix := getFirstDate(bugFixRate)
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
}

// graphOpenedVsClosedForUsers graphs many series:
// 1. Cumulative number of issues opened over time.
// 2. Cumulative number of issues closed over time.
// 3. Cumulative number of issues fixed over time for each user.
func graphOpenedVsClosedForUsers(issues []octokit.Issue, pulls []octokit.PullRequest, graphFileName string, users ...string) {
	// Get data for issues fixed for each user.
	var userSeries []series
	for _, userName := range users {
		newSeries := seriesFromMap(
			fmt.Sprintf("Total fixed by %s", userName),
			getBugFixRate(pulls, userName))
		userSeries = append(userSeries, newSeries)
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
}

// graphBugfixRateByUser graphs many series:
// 1. Cumulative number of issues opened over time.
// 2. Cumulative number of issues closed over time.
// 3. Cumulative number of issues fixed over time for each user who has
//    fixed at least a given number of issues.
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
}
