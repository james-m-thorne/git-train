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

func ConfigDeleteParent(currentBranch string) string {
	return fmt.Sprintf("git config --unset git-train.%s.parent 2> /dev/null", currentBranch)
}

func CheckoutNewBranch(branch string) string {
	return fmt.Sprintf("git checkout -b %s", branch)
}

func RebaseOntoParent(targetBranch string, parentBranch string, currentBranch string) string {
	return fmt.Sprintf("git rebase --onto %s %s %s", targetBranch, parentBranch, currentBranch)
}

func Delete(branch string) string {
	return fmt.Sprintf("git branch -D %s", branch)
}
