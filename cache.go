package main

import (
	"encoding/json"
	"os"

	"github.com/octokit/go-octokit/octokit"
)

// writeIssuesToCache serialises an array of issues and writes them to
// a json text file.
func writeIssuesToCache(fileName string, issues []octokit.Issue) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.Encode(issues)
}

// writeToCache serialises the array of pull requests and writes them
// to a json text file.
func writeToCache(fileName string, data []octokit.PullRequest) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.Encode(data)
}

// issuesFromCache reads an array of octokit issues from a json text
// file.
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

// pullsFromCache reads an array of pull requests from a json text
// file.
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

// getDataFromCache gets all issues and pull requests from the cache.
func getDataFromCache(issuesCache, pullsCache string) ([]octokit.Issue, []octokit.PullRequest) {
	return issuesFromCache(issuesCache), pullsFromCache(pullsCache)
}
