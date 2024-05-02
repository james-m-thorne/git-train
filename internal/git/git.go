package git

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
)

func ConfigGetAll() string {
	return "git config --list | grep git-train"
}

func ConfigGetMaster() string {
	return "git config --get git-train.master-branch || echo master"
}

func ConfigSetMaster(branch string) string {
	return fmt.Sprintf("git config git-train.master-branch %s", branch)
}

func GetCurrentBranch() string {
	return "git rev-parse --abbrev-ref HEAD"
}

func ConfigGetParent(branch string) string {
	return fmt.Sprintf("git config --get git-train.%s.parent", branch)
}

func ConfigSetParent(currentBranch string, newBranch string) string {
	return fmt.Sprintf("git config git-train.%s.parent %s", newBranch, currentBranch)
}

func CheckoutNewBranch(branch string) string {
	return fmt.Sprintf("git checkout -b %s", branch)
}

func RebaseOntoParent() (string, error) {
	currentBranch, err := command.GetOutput(GetCurrentBranch())
	if currentBranch == "" || err != nil {
		return "", fmt.Errorf("current branch not found")
	}
	parentBranch, err := command.GetOutput(ConfigGetParent(currentBranch))
	if parentBranch == "" || err != nil {
		return "", fmt.Errorf("no parent branch found for %s", currentBranch)
	}
	parentsParentBranch, err := command.GetOutput(ConfigGetParent(parentBranch))
	if parentsParentBranch == "" || err != nil {
		return "", fmt.Errorf("no parent branch found for %s", parentBranch)
	}
	return fmt.Sprintf("git rebase --onto %s %s %s", parentsParentBranch, parentBranch, currentBranch), nil
}
