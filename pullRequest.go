package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/octokit/go-octokit/octokit"
)

// list of keywords taken from https://help.github.com/articles/closing-issues-using-keywords/
const resolvesRegex = "[close | closes | closed | fix | fixes | fixed | resolve | resolves | resolved] #[0-9]+"

type pullRequest struct {
	pull             octokit.PullRequest
	referencedIssues []int
}

func newPull(base octokit.PullRequest) pullRequest {
	pull := pullRequest{
		pull: base,
	}
	rx := regexp.MustCompile(resolvesRegex)
	matches := rx.FindAllString(pull.pull.Body, -1)

	for _, st := range matches {
		// Get ID/Number of this issue
		index := strings.Index(st, "#")
		issueString := st[index+1:]
		issue, err := strconv.Atoi(issueString)
		if err != nil {
			fmt.Printf("Error getting issue resolved by pull request #%d: %v\n", base.Number, err)
		}
		pull.referencedIssues = append(pull.referencedIssues, issue)
	}
	return pull
}
