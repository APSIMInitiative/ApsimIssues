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

// Calculates the exact amount of time between two dates.
// https://stackoverflow.com/a/36531443
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// monthsBetween calculates the number of months between two dates.
func monthsBetween(a, b time.Time) int {
	_, month, _, _, _, _ := diff(a, b)
	return month
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

// getAllPullRequests uses the github api to fetch all pull requests
// (open and closed) on a given repo.
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

// Gets all data. Will attempt use the cache if the useCache global is
// set to true. Will call the github api otherwise.
func getData(client *octokit.Client) ([]octokit.Issue, []octokit.PullRequest) {
	// Only use cache if cache files are available.
	if useCache && fileExists(issuesCache) && fileExists(pullsCache) {
		fmt.Println("Fetching data from cache. This data is not live...")
		return getDataFromCache(issuesCache, pullsCache)
	}
	// Only show progress if not in quiet mode.
	issues, pulls := getDataFromGithub(client, owner, repo, !quiet)

	// Update cache for next time.
	writeToCache(pullsCache, pulls)
	writeIssuesToCache(issuesCache, issues)

	return issues, pulls
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
