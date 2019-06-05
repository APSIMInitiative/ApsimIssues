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

var quiet = false
var issuesCache = ".issues.cache.dat"
var pullsCache = ".pulls.cache.dat"
var useCache = true // Change to false for release!!!

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

	calcBugFixRate(username, client, "bugs.png")

	numIssuesByDate(client, owner, repo)
}
