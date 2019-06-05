package main

import (
	"bufio"
	"encoding/json"
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

// Reads an array of octokit issues from a json text file.
func issuesFromCache(fileName string) []octokit.Issue {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)

	// Read opening brace.
	_, err = decoder.Token()
	if err != nil {
		panic(err)
	}

	// Deserialise each value in the array.
	var issues []octokit.Issue
	for decoder.More() {
		var issue octokit.Issue
		err := decoder.Decode(&issue)
		if err != nil {
			panic(err)
		}
		issues = append(issues, issue)
	}
	return issues
}

// Reads an array of pull requests from a json text file.
func pullsFromCache(fileName string) []octokit.PullRequest {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)

	// Read opening brace.
	_, err = decoder.Token()
	if err != nil {
		panic(err)
	}

	// Deserialise each value in the array.
	var pulls []octokit.PullRequest
	for decoder.More() {
		var pull octokit.PullRequest
		err := decoder.Decode(&pull)
		if err != nil {
			panic(err)
		}
		pulls = append(pulls, pull)
	}
	return pulls
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

func readFromCache(fileName string, object *interface{}) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	decoder.Decode(object)
}

// Serialises the array of issues and writes them to a cache file.
func writeIssuesToCache(fileName string, issues []octokit.Issue) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.Encode(issues)
}

// Serialises the array of issues and writes them to a cache file.
func writeToCache(fileName string, data []octokit.PullRequest) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.Encode(data)
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
