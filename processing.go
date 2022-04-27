package main

import (
	"time"

	"github.com/octokit/go-octokit/octokit"
)

// getNumOpenIssues takes an array of issues and returns the number of
// issues which are open.
func getNumOpenIssues(issues []octokit.Issue) int {
	var sum int
	for _, issue := range issues {
		if issue.ClosedAt == nil {
			sum++
		}
	}
	return sum
}

// getNumClosedIssues takes an array of issues and returns the number of
// issues which are closed.
func getNumClosedIssues(issues []octokit.Issue) int {
	var sum int
	for _, issue := range issues {
		if issue.ClosedAt != nil {
			sum++
		}
	}
	return sum
}

// getNumOpenPullRequests takes an array of pull requests and returns
// the number of pull requests which are open.
func getNumOpenPullRequests(pulls []octokit.PullRequest) int {
	var sum int
	for _, pull := range pulls {
		if pull.ClosedAt == nil {
			sum++
		}
	}
	return sum
}

// getNumClosedPullRequests takes an array of pull requests and returns
// the number of pull requests which are open.
func getNumClosedPullRequests(pulls []octokit.PullRequest) int {
	var sum int
	for _, pull := range pulls {
		if pull.ClosedAt != nil {
			sum++
		}
	}
	return sum
}

func getNumIssuesOpenedBy(issues []octokit.Issue, user string) int {
	var sum int
	for _, issue := range issues {
		if issue.User.Login == user {
			sum++
		}
	}
	return sum
}

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
			addAfterDate(&issues, *pull.pull.ClosedAt, len(pull.referencedIssues))
		}
	}
	return issues
}

// issuesFixedByStaleBot returns all issues which were fixed by StaleBot.
func issuesFixedByStaleBot(issues []octokit.Issue) []octokit.Issue {
	closedIssues := filterIssues(issues, func(i octokit.Issue) bool {
		return i.ClosedAt != nil
	})
	return filterIssues(closedIssues, func(i octokit.Issue) bool {
		return issueHasLabel(i, "stale")
	})
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
// issues on that date.
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

// isBug checks if an issue is a bug
// todo - refactor
func isBug(issue octokit.Issue) bool {
	var labels []string
	for _, label := range issue.Labels {
		labels = append(labels, label.Name)
	}

	return indexOfString(labels, "bug") >= 0
}

// isOpen checks if an issue is open
func isOpen(issue octokit.Issue) bool {
	return issue.ClosedAt == nil
}

// isClosed checks if an issue is closed
func isClosed(issue octokit.Issue) bool {
	return !isOpen(issue)
}

// hasLabel checks if an issue has a given label.
func hasLabel(issue octokit.Issue, label string) bool {
	for _, label := range issue.Labels {
		if label.Name == settings.LabelFilter {
			return true
		}
	}
	return false
}

// hasLabel checks if an issue has a given label.
func issueHasLabel(issue octokit.Issue, label string) bool {
	for _, lbl := range issue.Labels {
		if lbl.Name == label {
			return true
		}
	}
	return false
}

// issuesWithLabel takes a list of issues and returns those issues with a given label.
func issuesWithLabel(issues []octokit.Issue, label string) []octokit.Issue {
	return filterIssues(issues, func(issue octokit.Issue) bool {
		return hasLabel(issue, label)
	})
}

func getIssueWithID(issues []octokit.Issue, id int) *octokit.Issue {
	for _, issue := range issues {
		if issue.Number == id {
			return &issue
		}
	}
	return nil
}

func pullsWithLabel(pulls []octokit.PullRequest, issues []octokit.Issue, label string) []octokit.PullRequest {
	return filterPullRequests(pulls, func(pull octokit.PullRequest) bool {
		pullRequest := newPull(pull)
		for _, issueID := range pullRequest.referencedIssues {
			issue := getIssueWithID(issues, issueID)
			if issue != nil && hasLabel(*issue, label) {
				return true
			}
		}
		return false
	})
}

// bugsFixedSince returns the number of bugs (issues with the label 'bug') fixed since a given date
func bugsFixedSince(issues []octokit.Issue, date time.Time) int {
	return len(filterIssues(issues, func(issue octokit.Issue) bool {
		return isBug(issue) && issue.ClosedAt != nil && (issue.ClosedAt.After(date) || sameDay(*issue.ClosedAt, date))
	}))
}

// issuesFixedSince returns the number of issues fixed since a given date
func issuesFixedSince(issues []octokit.Issue, date time.Time) int {
	return len(filterIssues(issues, func(issue octokit.Issue) bool {
		return issue.ClosedAt != nil && (issue.ClosedAt.After(date) || sameDay(*issue.ClosedAt, date))
	}))
}

// get all issues grouped by the user who created the issues.
func getIssuesGroupedByAuthor(issues []octokit.Issue) map[string][]octokit.Issue {
	groups := make(map[string][]octokit.Issue)
	for _, issue := range issues {
		groups[issue.User.Login] = append(groups[issue.User.Login], issue)
	}
	return groups
}
