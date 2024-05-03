package git

import (
	"fmt"
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
	return "git branch --show-current"
}

func ConfigGetParent(branch string) string {
	return fmt.Sprintf("git config --get git-train.%s.parent", branch)
}

func ConfigGetChild(branch string) string {
	return fmt.Sprintf("git config --list | grep .parent=%s", branch)
}

func ConfigSetParent(currentBranch string, parentBranch string) string {
	return fmt.Sprintf("git config git-train.%s.parent %s", currentBranch, parentBranch)
}

func ConfigDeleteParent(currentBranch string) string {
	return fmt.Sprintf("git config --unset git-train.%s.parent 2> /dev/null", currentBranch)
}

func CheckoutNewBranch(branch string) string {
	return fmt.Sprintf("git checkout -b %s", branch)
}

func Checkout(branch string) string {
	return fmt.Sprintf("git checkout %s", branch)
}

func Rebase(targetBranch string) string {
	return fmt.Sprintf("git rebase %s", targetBranch)
}

func RebaseOntoTarget(targetBranch string, ignoreBranch string, currentBranch string) string {
	return fmt.Sprintf("git rebase --onto %s %s %s", targetBranch, ignoreBranch, currentBranch)
}

func Delete(branch string) string {
	return fmt.Sprintf("git branch -D %s", branch)
}

func Push() string {
	return "git push"
}
