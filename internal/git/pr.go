package git

import (
	"encoding/json"
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"regexp"
	"strings"
)

type PullRequestState struct {
	HeadRefName string `json:"headRefName"`
	Body        string `json:"body"`
	Number      int    `json:"number"`
}

func GetBranchStackPullRequests(branchStack []string) map[string]PullRequestState {
	prsJson := command.GetOutputFatal(GitHubPrListBranchStack(branchStack))

	var pullRequests []PullRequestState
	err := json.Unmarshal([]byte(prsJson), &pullRequests)
	if err != nil {
		command.PrintFatalError("Error parsing pull requests: %s", err)
	}

	branchPullRequest := map[string]PullRequestState{}
	for _, branch := range branchStack {
		for _, pr := range pullRequests {
			if branch == pr.HeadRefName {
				branchPullRequest[branch] = pr
			}
		}
	}

	return branchPullRequest
}

func UpdatePullRequestBodies(branchStack []string, branchPullRequests map[string]PullRequestState) map[string]PullRequestState {
	prList := ""
	for i, branch := range branchStack {
		if pr, ok := branchPullRequests[branch]; ok {
			prList += fmt.Sprintf("%d. #%d\n", i+1, pr.Number)
		} else {
			prList += fmt.Sprintf("%d. #%s\n", i+1, branch)
		}
	}

	body := fmt.Sprintf(`
<!---GitTrainStart--->
**Pull Request Stack:**

%s

Managed with ❤️ by [james-m-thorne/git-train](https://github.com/james-m-thorne/git-train)
<!---GitTrainEnd--->
`, prList)

	// Compile the regular expression to match everything between
	// <!---GitTrainStart---> and <!---GitTrainEnd--->
	re := regexp.MustCompile(`(?s)<!---GitTrainStart--->.*?<!---GitTrainEnd--->`)

	for _, branch := range branchStack {
		if pr, ok := branchPullRequests[branch]; ok {
			if re.MatchString(pr.Body) {
				pr.Body = re.ReplaceAllString(pr.Body, body)
			} else {
				pr.Body = fmt.Sprintf("%s\n\n%s", pr.Body, body)
			}
			branchPullRequests[branch] = pr
		}
	}

	return branchPullRequests
}

func GitHubPrListBranchStack(branchStack []string) string {
	branchesWithQuotes := make([]string, len(branchStack))
	for i, branch := range branchStack {
		branchesWithQuotes[i] = fmt.Sprintf("\"%s\"", branch)
	}
	return fmt.Sprintf("gh pr list --author \"@me\" --json number,body,headRefName --jq '[.[] | select(.headRefName | IN(%s))]'", strings.Join(branchesWithQuotes, ","))
}
