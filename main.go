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
var cache = ".cache.txt"

func main() {
	// We expect the user to pass in a username as a command line argument.
	var username string
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No username received as command line argument. Defaulting to hol430...")
		username = "hol430"
	} else {
		username = args[1]
	}
	if len(args) > 2 {
		if args[2] == "-q" {
			quiet = true
		}
	}
	fmt.Printf("username=%s\n", username) // Temp code to fix unused var warning!!!! Don't commit!!

	auth := getAuth("credentials.dat")
	client := octokit.NewClient(auth)

	pulls := pullsByUser(username, client)
	graphFile := "bugs.png"
	createBugfixGraph(pulls, graphFile)
	fmt.Printf("Generated graph '%v'\n", graphFile)

	dataFile := "issues.csv"
	exportToCsv(dataFile, pulls)
	fmt.Printf("Generating data file '%v'\n", dataFile)
	fmt.Printf("%s has resolved %d issues.\n", username, numIssuedResolved(pulls))

	numIssuesByDate(client, owner, repo)
}
