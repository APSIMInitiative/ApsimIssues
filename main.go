package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/octokit/go-octokit/octokit"
)

const (
	owner       = "APSIMInitiative"
	repo        = "ApsimX"
	issuesCache = ".issues.cache"
	pullsCache  = ".pulls.cache"
)

var (
	settings options
)

func main() {
	args, err := flags.Parse(&settings)
	if err != nil {
		if flags.WroteHelp(err) {
			return
		}
		panic(err)
	}
	if len(args) > 0 {
		// If there are any leftover unrecognised arguments, throw a fatal
		panic(fmt.Sprintf("Error: unrecognised arguments: %v", args))
	}

	auth := getAuth("credentials.dat")
	client := octokit.NewClient(auth)
	issues, pullRequests := getData(client)

	if settings.LabelFilter != "" {
		if !settings.Quiet {
			fmt.Printf("Filtering on issues with the %s label...\n", settings.LabelFilter)
		}
		issues = issuesWithLabel(issues, settings.LabelFilter)
		pullRequests = pullsWithLabel(pullRequests, issues, settings.LabelFilter)
	}

	// Diagnostics
	if !settings.Quiet {
		fmt.Printf("Owner:                                  %s\n", owner)
		fmt.Printf("Repo:                                   %s\n", repo)
		fmt.Printf("User:                                   %s\n\n", settings.Username)
	}

	fmt.Printf("Number of open issues:                      %d\n", getNumOpenIssues(issues))
	fmt.Printf("Number of closed issues:                    %d\n", getNumClosedIssues(issues))
	fmt.Printf("Number of open pull requests:               %d\n", getNumOpenPullRequests(pullRequests))
	fmt.Printf("Number of closed pull requests:             %d\n", getNumClosedPullRequests(pullRequests))
	fmt.Printf("Number of issues opened by %s:              %d\n", settings.Username, getNumIssuesOpenedBy(issues, settings.Username))

	since := settings.Since()
	issues = filterIssues(issues, func(issue octokit.Issue) bool {
		return issue.CreatedAt.After(since)
	})
	pullRequests = filterPullRequests(pullRequests, func(pull octokit.PullRequest) bool {
		return pull.MergedAt != nil && pull.MergedAt.After(since)
	})
	fmt.Printf("Number of bugs closed since %s:             %d\n", since.Format("2/1/2006"), bugsFixedSince(issues, since))
	fmt.Printf("Number of issues closed since %s:           %d\n\n", since.Format("2/1/2006"), issuesFixedSince(issues, since))

	// Graphs
	graphBugFixRate(pullRequests, settings.Username, "bugs.png")
	graphIssuesByDate(issues, "openIssues.png")
	graphOpenedVsClosed(issues, "openedVsClosed.png")
    graphOpenedVsClosedForUser(issues, pullRequests, settings.Username, "closedByUser.png")
	graphOpenedVsClosedForUsers(issues, pullRequests, "fixersComparison.png", settings.Username, "zur003", "hol353")
	graphBugfixRateByUser(issues, pullRequests, "fixersComparisonByBugCount.png", 100)
	graphBugfixRateByUser(issues, pullRequests, "allfixersComparison.png", -1)
	graphIssuesOpenedByUser(issues, 50, "issuesOpenedByUser.png")
}
