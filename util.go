package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/octokit/go-octokit/octokit"
)

// Reads `filename` and returns an authentication method
func getAuth(filename string) octokit.AuthMethod {
	credentials, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}
	var username, password string
	scanner := bufio.NewScanner(strings.NewReader(string(credentials)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "username=") {
			username = strings.TrimPrefix(line, "username=")
		}
		if strings.HasPrefix(line, "password=") {
			password = strings.TrimPrefix(line, "password=")
		}
		if strings.HasPrefix(line, "token=") {
			token := strings.TrimPrefix(line, "token=")
			return octokit.TokenAuth{AccessToken: token}
		}
	}

	return octokit.BasicAuth{Login: username, Password: password}
}

func sortKeys(m map[time.Time]int) []time.Time {
	keys := make([]time.Time, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
	return keys
}

// Gets all issues (open and closed) in a repository.
func getAllIssues(client *octokit.Client, owner, repo string, showProgress bool) (issues []octokit.Issue) {
	apsimURL := octokit.Hyperlink("repos/{owner}/{repo}/issues?state={state}")
	first := true
	var numIssues int
	for &apsimURL != nil {
		issuesSubset, result := client.Issues().All(&apsimURL, octokit.M{
			"owner": owner,
			"repo":  repo,
			"state": "all",
		})
		if result.HasError() {
			panic(result)
		}
		for _, issue := range issuesSubset {
			if first && showProgress {
				fmt.Printf("Updating numIssues to %d\n", issue.Number)
				numIssues = issue.Number
				first = false
			}
			if showProgress {
				percentDone := 100.0 * float64(numIssues-issue.Number) / float64(numIssues)
				fmt.Printf("Fetching issues: %.2f%%...\r", percentDone)
			}
			if issue.PullRequest.HTMLURL == "" {
				issues = append(issues, issue)
			}
		}
		if result.NextPage == nil {
			break
		}
		apsimURL = *result.NextPage
	}
	if showProgress {
		fmt.Printf("Fetching issues: 100.00%%...\n")
	}
	return
}

func getAllPullRequests(client *octokit.Client, owner, repo string, showProgress bool) []octokit.PullRequest {
	var pulls []octokit.PullRequest
	apsimURL := octokit.Hyperlink("repos/{owner}/{repo}/pulls?state=closed")

	first := true
	var numPullRequests int
	var percentDone float64
	for &apsimURL != nil {
		url, err := apsimURL.Expand(octokit.M{
			"owner": owner,
			"repo":  repo,
		})
		if err != nil {
			panic(err)
		}

		allPulls, result := client.PullRequests(url).All()
		if result.HasError() {
			panic(result)
		}
		for _, pull := range allPulls {
			if first {
				numPullRequests = pull.Number
				first = false
			}

			if showProgress {
				percentDone = 100.0 * float64(numPullRequests-pull.Number) / float64(numPullRequests)
				fmt.Printf("Fetching pull requests: %.2f%%...\r", percentDone)
			}

			pulls = append(pulls, pull)
		}
		if result.NextPage == nil {
			break
		}
		apsimURL = *result.NextPage
	}
	if showProgress {
		fmt.Printf("Fetching pull requests: 100.00%%...\r")
	}
	return pulls
}

// Gets all issues (open and closed) in a repository.
func getDataFromGithub(client *octokit.Client, owner, repo string, showProgress bool) (issues []octokit.Issue, pulls []octokit.PullRequest) {
	// TODO : combine these methods.
	issues = getAllIssues(client, owner, repo, showProgress)
	pulls = getAllPullRequests(client, owner, repo, showProgress)
	return
}

// Checks if a file exists
func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		panic(fmt.Sprintf("Error: Unable to determine if file '%s' exists", fileName))
	}
}
