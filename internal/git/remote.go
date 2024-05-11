package git

import (
	"fmt"
)

func GitHubPrCreate(parentBranch string) string {
	return fmt.Sprintf("gh pr create --base %s --web", parentBranch)
}

func GitHubPrEditBody(number int, body string) string {
	return fmt.Sprintf("gh pr edit %d --body-file - <<EOF\n%s\nEOF", number, body)
}

func GitHubPrState() string {
	return "gh pr status --json state --jq '.currentBranch.state'"
}
