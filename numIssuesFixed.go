package main
import (
	"fmt"
	"github.com/octokit/go-octokit/octokit"
	"os"
	"regexp"
	"io/ioutil"
	"bufio"
	"strings"
	"log"
)
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
	
	apsimUrl := octokit.Hyperlink("repos/APSIMInitiative/ApsimX/pulls?state=closed")
	auth := octokit.BasicAuth{Login: apiUsername, Password: password}
	client := octokit.NewClient(auth)
	// list of keywords taken from https://help.github.com/articles/closing-issues-using-keywords/
	resolvesRegex := regexp.MustCompile("[close | closes | closed | fix | fixes | fixed | resolve | resolves | resolved] #[0-9]")
	numResolvedIssues := 0
	first := true
	var numPullRequests int
	var percentDone float64
	
	for &apsimUrl != nil {
		url, err := apsimUrl.Expand(nil)
		if err != nil {
			log.Fatal(err)
		}
		
		pulls, result := client.PullRequests(url).All()
		if result.HasError() {
			fmt.Println("error:", result)
			return
		}
		for _, pull := range pulls {
			if first {
				numPullRequests = pull.Number
				first = false
			}
			percentDone = 100.0 * float64(numPullRequests - pull.Number) / float64(numPullRequests)
			fmt.Printf("Working...%.2f%%\r", percentDone)
			if pull.User.Login == username {
				numResolvedIssues += len(resolvesRegex.FindAllStringIndex(pull.Body, -1))
			}
		}
		if result.NextPage == nil {
			break
		}
		apsimUrl = *result.NextPage
	}
	fmt.Printf("%s has resolved %d issues.\n", username, numResolvedIssues)
	fmt.Println("Press enter to exit.")
	fmt.Scanln()
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