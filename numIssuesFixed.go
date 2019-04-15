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

func pullsByUser(username string, client *octokit.Client) []pullRequest {
	var pulls []pullRequest
	apsimURL := octokit.Hyperlink("repos/APSIMInitiative/ApsimX/pulls?state=closed")

	first := true
	var numPullRequests int
	var percentDone float64
	for &apsimURL != nil {
		url, err := apsimURL.Expand(nil)
		if err != nil {
			log.Fatal(err)
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
			percentDone = 100.0 * float64(numPullRequests-pull.Number) / float64(numPullRequests)
			fmt.Printf("Working...%.2f%%\r", percentDone)
			if pull.User.Login == username {
				pulls = append(pulls, newPull(pull))
			}
		}
		if result.NextPage == nil {
			break
		}
		apsimURL = *result.NextPage
	}
	return pulls
}

func numIssuedResolved(pulls []pullRequest) int {
	n := 0
	for _, pull := range pulls {
		n += len(pull.referencedIssues)
	}
	return n
}

func getIssuesByDate(pulls []pullRequest) map[time.Time]int {
	issues := make(map[time.Time]int)
	for _, pull := range pulls {
		issues[*pull.pull.ClosedAt] += len(pull.referencedIssues)
	}
	return issues
}

// Exports a csv file containing two columns: date and num Issued
// resolved on that date
func exportToCsv(filename string, pulls []pullRequest) {
	if _, err := os.Stat("/path/to/whatever"); err == nil {
		// Delete file if it exists
		err := os.Remove(filename)
		if err != nil {
			panic(fmt.Sprintf("Unable to delete file %v: %v", filename, err))
		}
	}

	// Open the file in append mode
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	issuesBydate := getIssuesByDate(pulls)
	sortedDates := sortKeys(issuesBydate)
	for _, date := range sortedDates {
		numIssues := issuesBydate[date]
		str := fmt.Sprintf("%v,%d\n", date, numIssues)
		_, err := f.Write([]byte(str))
		if err != nil {
			panic(err)
		}
	}
}

func sortKeys(m map[time.Time]int) ([]time.Time) {
	keys := make([]time.Time, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j])})
	return keys
}
func main() {
	// We expect the user to pass in a username as a command line argument.
	var username string
	args := os.Args
	if len(args) != 2 {
		fmt.Println("No username received as command line argument. Defaulting to hol430...")
		username = "hol430"
	} else {
		username = args[1]
	}

	apiUsername, password := getCredentials()

	auth := octokit.BasicAuth{Login: apiUsername, Password: password}
	client := octokit.NewClient(auth)
	pulls := pullsByUser(username, client)
	exportToCsv("issues.dat", pulls)
	fmt.Printf("%s has resolved %d issues.\n", username, numIssuedResolved(pulls))
}

func getCredentials() (username string, password string) {
	credentials, err := ioutil.ReadFile("credentials.dat")

	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(credentials)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "username=") {
			username = strings.TrimPrefix(line, "username=")
		}
		if strings.HasPrefix(line, "password=") {
			password = strings.TrimPrefix(line, "password=")
		}
	}
	return
}
