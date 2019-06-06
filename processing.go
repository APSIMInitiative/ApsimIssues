package main

import (
	"fmt"
	"time"

	"github.com/octokit/go-octokit/octokit"
)

// pullsByUser takes an array of pull requests and a username, and
// returns all pull requests created by that user.
func pullsByUser(username string, allPulls []octokit.PullRequest) []pullRequest {
	var pulls []pullRequest
	for _, pull := range allPulls {
		if pull.User.Login == username {
			pulls = append(pulls, newPull(pull))
		}
	}
	return pulls
}

// pullsGroupedByUser takes an array of pull requests and returns a map
// of usernames to an array of pull requests created by that user.
func pullsGroupedByUser(allPulls []octokit.PullRequest) (result map[string][]pullRequest) {
	result = make(map[string][]pullRequest)
	for _, pull := range allPulls {
		result[pull.User.Login] = append(result[pull.User.Login], newPull(pull))
	}
	return
}

// getCumIssuesByDate takes a list of pull requests and returns a map
// of dates to the number of issues resolved on and before that date.
func getCumIssuesByDate(pulls []pullRequest) map[time.Time]int {
	issues := make(map[time.Time]int)
	for _, pull := range pulls {
		if pull.pull.ClosedAt != nil {
			issues[*pull.pull.ClosedAt] = 0
		}
	}
	for _, pull := range pulls {
		if pull.pull.ClosedAt != nil {
			incrementAfterDate(&issues, *pull.pull.ClosedAt)
		}
	}
	return issues
}

// getIssuesByDate takes an array of pull requests and returns a map of
// dates to the number of issues resolved on that date.
func getIssuesByDate(pulls []pullRequest) map[time.Time]int {
	issues := make(map[time.Time]int)
	for _, pull := range pulls {
		if pull.pull.ClosedAt != nil {
			issues[*pull.pull.ClosedAt] += len(pull.referencedIssues)
		}
	}
	return issues
}

// getBugFixRate gets the number of bugs fixed over time by a given
// user. Returns a map of dates to number of bugs fixed on that date.
func getBugFixRate(allPulls []octokit.PullRequest, username string) map[time.Time]int {
	// Filter out pull requests not created by the user.
	pulls := pullsByUser(username, allPulls)

	// Generate a map of dates to cumulative number of issues
	// referenced in pull requests.
	return getCumIssuesByDate(pulls)
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

// getCumIssuesClosedByDate gets a map of dates to the number of closed
// issuse on that date.
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
	fmt.Printf("Generated graph '%s'\n", graphFileName)
	if bugFixRate != nil {
		fmt.Printf("%s has resolved %d issues.\n", username, bugFixRate[getLastDate(bugFixRate)])
	}
}
