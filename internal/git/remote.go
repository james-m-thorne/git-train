package git

import "fmt"

func GitHubPrCreate(parentBranch string) string {
	return fmt.Sprintf("gh pr create --base %s --web", parentBranch)
}

func GitHubPrState() string {
	return "gh pr status --json state --jq '.currentBranch.state'"
}

func GitHubPrView() string {
	return "gh pr view --web"
}

func GitHubPrBody() string {
	return "gh pr view --json body --jq '.body'"
}
