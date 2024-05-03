package git

import "github.com/james-m-thorne/git-train/internal/command"

func GetBranchStack(currentBranch string, includeMaster bool) []string {
	var branchStack []string
	masterBranch := ""
	if !includeMaster {
		masterBranch, _ = command.GetOutput(ConfigGetMaster())
	}
	for currentBranch != masterBranch {
		branchStack = append(branchStack, currentBranch)
		parentBranch, err := command.GetOutput(ConfigGetParent(currentBranch))
		if err != nil {
			break
		}
		currentBranch = parentBranch
	}
	return branchStack
}
