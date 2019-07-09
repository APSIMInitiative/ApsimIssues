package main

import (
	"fmt"
	"os"

	"github.com/octokit/go-octokit/octokit"
)

const (
	owner = "APSIMInitiative"
	repo  = "ApsimX"
)

var (
	quiet       = false
	issuesCache = ".issues.cache"
	pullsCache  = ".pulls.cache"
	useCache    = true // Change to false for release!!!
)

func main() {
	// We expect the user to pass in a username as a command line argument.
	var username string
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No username received as command line argument. Defaulting to hol430...")
		username = "hol430"
	} else {
		username = args[1]
		fmt.Printf("username=%s\n", username)
	}
	if len(args) > 2 {
		if args[2] == "-q" {
			quiet = true
		}
	}

	auth := getAuth("credentials.dat")
	client := octokit.NewClient(auth)
	issues, pullRequests := getData(client)

	// Diagnostics
	if !quiet {
		fmt.Printf("Owner:                      %s\n", owner)
		fmt.Printf("Repo:                       %s\n", repo)
		fmt.Printf("User:                       %s\n\n", username)
	}
	fmt.Printf("Number of open issues:          %d\n", getNumOpenIssues(issues))
	fmt.Printf("Number of closed issues:        %d\n", getNumClosedIssues(issues))
	fmt.Printf("Number of open pull requests:   %d\n", getNumOpenPullRequests(pullRequests))
	fmt.Printf("Number of closed pull requests: %d\n\n", getNumClosedPullRequests(pullRequests))

	// Graphs
	graphBugFixRate(pullRequests, username, "bugs.png")
	graphIssuesByDate(issues, "openIssues.png")
	graphOpenedVsClosed(issues, "openedVsClosed.png")
	graphOpenedVsClosedForUser(issues, pullRequests, username, "closedByUser.png")
	graphOpenedVsClosedForUsers(issues, pullRequests, "fixersComparison.png", username, "zur003", "hol353")
	graphBugfixRateByUser(issues, pullRequests, "fixersComparison.png", 100)
	graphBugfixRateByUser(issues, pullRequests, "allfixersComparison.png", -1)

}
