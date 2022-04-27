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

// getAuth reads a file and returns a github authentication method.
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

// sortKeys returns a slice of all keys in a map, sorted in
// chronological ascending order.
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

// indexOf searches a slice for a certain item and returns its index,
// or -1 if not found.
func indexOfString(arr []string, item string) int {
	for i, str := range arr {
		if str == item {
			return i
		}
	}
	return -1
}

// indexOf searches a slice of dates for a certain date. Returns the
// index of the item in the slice, or -1 if not found.
func indexOf(dates []time.Time, date time.Time) int {
	for i, dt := range dates {
		if sameDay(date, dt) {
			return i
		}
	}
	return -1
}

// diff calculates the exact amount of time between two dates.
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

// sameDay returns true iff two time objects represent the same day.
func sameDay(d1, d2 time.Time) bool {
	return d1.Year() == d2.Year() && d1.YearDay() == d2.YearDay()
}

// monthsBetween calculates the number of months between two dates.
func monthsBetween(a, b time.Time) int {
	_, month, _, _, _, _ := diff(a, b)
	return month
}

// incrementAfterDate increments all values in the map whose date key
// lies on or after a given date.
func incrementAfterDate(issues *map[time.Time]int, date time.Time) {
	for key := range *issues {
		if key.After(date) || key == date {
			(*issues)[key]++
		}
	}
}

// addAfterDate adds a given value to all values in the map whose date
// key lies on or after a given date.
func addAfterDate(issues *map[time.Time]int, date time.Time, value int) {
	for key := range *issues {
		if key.After(date) || key == date {
			(*issues)[key] += value
		}
	}
}

// decrementAfterDate decrements all values in the map whose date key
// lies on or after a given date.
func decrementAfterDate(issues *map[time.Time]int, date time.Time) {
	for key := range *issues {
		if key.After(date) || key == date {
			(*issues)[key]--
		}
	}
}

// getLastDate finds the (chronologically) last date in a map.
func getLastDate(data map[time.Time]int) time.Time {
	var lastDate time.Time
	first := true
	for date := range data {
		if first {
			lastDate = date
			first = false
		} else if date.After(lastDate) {
			lastDate = date
		}
	}
	return lastDate
}

// getFirstDate finds the (chronologically) first date in a map.
func getFirstDate(data map[time.Time]int) time.Time {
	var firstDate time.Time
	first := true
	for date := range data {
		if first {
			firstDate = date
			first = false
		} else if date.Before(firstDate) {
			firstDate = date
		}
	}
	return firstDate
}

// getAllIssues gets all issues (open and closed) on a github
// repository.
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
				fmt.Printf("\rFetching issues: %.2f%%...", percentDone)
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
		fmt.Printf("\rFetching issues: 100.00%%...\n")
	}
	return
}

// getAllPullRequests gets all pull requests (open and closed) on a
// github repository.
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
				fmt.Printf("\rFetching pull requests: %.2f%%...", percentDone)
			}

			pulls = append(pulls, pull)
		}
		if result.NextPage == nil {
			break
		}
		apsimURL = *result.NextPage
	}
	if showProgress {
		fmt.Printf("\rFetching pull requests: 100.00%%...\n")
	}
	return pulls
}

// getDataFromGithub gets all issues and pull requests on a github
// repository by calling the github API.
func getDataFromGithub(client *octokit.Client, owner, repo string, showProgress bool) (issues []octokit.Issue, pulls []octokit.PullRequest) {
	// TODO : combine these methods.
	issues = getAllIssues(client, owner, repo, showProgress)
	pulls = getAllPullRequests(client, owner, repo, showProgress)
	return
}

// getData gets all data. Will attempt use the cache if the useCache
// global is set to true. Will get the data from github otherwise.
func getData(client *octokit.Client) ([]octokit.Issue, []octokit.PullRequest) {
	// Only use cache if cache files are available.
	if settings.UseCache && fileExists(issuesCache) && fileExists(pullsCache) {
		fmt.Println("Fetching data from cache. This data is not live...")
		return getDataFromCache(issuesCache, pullsCache)
	}
	// Only show progress if not in quiet mode.
	issues, pulls := getDataFromGithub(client, owner, repo, !settings.Quiet)

	// Update cache for next time.
	writeToCache(pullsCache, pulls)
	writeIssuesToCache(issuesCache, issues)

	return issues, pulls
}

// fileExists checks if a file exists
func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		panic(fmt.Sprintf("Error: Unable to determine if file '%s' exists", fileName))
	}
}

// filterIssues returns a deep clone of a slice of issues, filtered on
// a given predicate.
func filterIssues(issues []octokit.Issue, condition func(octokit.Issue) bool) []octokit.Issue {
	var result []octokit.Issue

	for _, issue := range issues {
		if condition(issue) {
			result = append(result, issue)
		}
	}

	return result
}

// filterPullRequests returns a deep clone of a slice of pull requests,
// filtered on a given predicate.
func filterPullRequests(pulls []octokit.PullRequest, condition func(octokit.PullRequest) bool) []octokit.PullRequest {
	var result []octokit.PullRequest

	for _, pull := range pulls {
		if condition(pull) {
			result = append(result, pull)
		}
	}

	return result
}

// filterIssueGropu returns a deep clone of a map of strings to an array of issues.
// The returned object contains only those key/value pairs which satisfy a condition.
func filterIssueGroup(issues map[string][]octokit.Issue, condition func([]octokit.Issue) bool) map[string][]octokit.Issue {
	result := make(map[string][]octokit.Issue)

	for key, val := range issues {
		if condition(val) {
			result[key] = val
		}
	}

	return result
}

// filterIssueGropu returns a deep clone of a map of strings to an array of issues.
// The returned object contains only those key/value pairs which satisfy a condition.
func filterIssueGroupIssues(issues map[string][]octokit.Issue, condition func(octokit.Issue) bool) map[string][]octokit.Issue {
	result := make(map[string][]octokit.Issue)

	for key, val := range issues {
		for _, issue := range val {
			if condition(issue) {
				result[key] = append(result[key], issue)
			}
		}
	}

	return result
}
