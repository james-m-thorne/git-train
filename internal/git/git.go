package git

import "fmt"

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

func ConfigSetParent(currentBranch string, newBranch string) string {
	return fmt.Sprintf("git config git-train.%s.parent %s", newBranch, currentBranch)
}

func CheckoutNewBranch(branch string) string {
	return fmt.Sprintf("git checkout -b %s", branch)
}
