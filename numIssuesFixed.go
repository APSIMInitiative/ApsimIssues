package main

import (
	"fmt"
	"time"

	"github.com/octokit/go-octokit/octokit"
)

func pullsByUser(username string, allPulls []octokit.PullRequest) []pullRequest {
	var pulls []pullRequest
	for _, pull := range allPulls {
		if pull.User.Login == username {
			pulls = append(pulls, newPull(pull))
		}
	}
	return pulls
}

func numIssuedResolved(pulls []pullRequest) int {
	n := 0
	for _, pull := range pulls {
		n += len(pull.referencedIssues)
	}
	return n
}

func getIssuesByDate(pulls []pullRequest) map[time.Time]int {
	issues := make(map[time.Time]int)
	for _, pull := range pulls {
		if pull.pull.ClosedAt != nil {
			issues[*pull.pull.ClosedAt] += len(pull.referencedIssues)
		}
	}
	return issues
}

func graphBugFixRate(allPulls []octokit.PullRequest, username, graphFileName string) {
	// Filter out pull requests not created by the user.
	pulls := pullsByUser(username, allPulls)

	// Generate a map of dates to number of issues referenced in pull requests.
	issuesByDate := getIssuesByDate(pulls)

	// Convert the map into two arrays. Maps have no concept of order,
	// but the data needs to be ordered by date ascending before we can
	// graph it.
	dates := sortKeys(issuesByDate)
	var issues []int // unused
	var cumIssues []int
	sum := 0
	for _, date := range dates {
		issues = append(issues, issuesByDate[date])
		sum += issuesByDate[date]
		cumIssues = append(cumIssues, sum)
	}

	createLinePlot(
		dates,
		cumIssues,
		fmt.Sprintf("Cumulative bugs fixed over time by %s", username),
		"Date",
		"Total Number of Issues Resolved",
		graphFileName)
	fmt.Printf("Generated graph '%s'\n", graphFileName)
	fmt.Printf("%s has resolved %d issues.\n", username, numIssuedResolved(pulls))

}
