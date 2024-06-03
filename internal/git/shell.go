package git

import (
	"fmt"
	"strings"
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

func ConfigGetRemote() string {
	return "git config --get git-train.remote || echo origin"
}

func ConfigSetRemote(remote string) string {
	return fmt.Sprintf("git config git-train.remote %s", remote)
}

func ConfigGetMergedCompletedBranches(mergedBranch string) string {
	return fmt.Sprintf("git config --get git-train.set-merged.%s.completed-branches || echo ''", mergedBranch)
}

func ConfigSetMergedCompletedBranches(mergedBranch string, branches []string) string {
	return fmt.Sprintf("git config git-train.set-merged.%s.completed-branches %s", mergedBranch, strings.Join(branches, ","))
}

func ConfigDeleteMergedCompletedBranches(mergedBranch string) string {
	return fmt.Sprintf("git config --unset git-train.set-merged.%s.completed-branches", mergedBranch)
}

func ConfigGetSyncCompletedBranches() string {
	return "git config --get git-train.sync.completed-branches || echo ''"
}

func ConfigSetSyncCompletedBranches(branches []string) string {
	return fmt.Sprintf("git config git-train.sync.completed-branches %s", strings.Join(branches, ","))
}

func ConfigDeleteSyncCompletedBranches() string {
	return "git config --unset git-train.sync.completed-branches"
}

func ShowCurrentBranch() string {
	return "git branch --show-current"
}

func CheckBranchExists(branch string) string {
	return fmt.Sprintf("git rev-parse --verify %s", branch)
}

func ConfigGetParent(branch string) string {
	return fmt.Sprintf("git config --get git-train.%s.parent", branch)
}

func ConfigGetChild(branch string) string {
	return fmt.Sprintf("git config --list | grep '.parent=%s$'", branch)
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

func Merge(targetBranch string) string {
	return fmt.Sprintf("git merge %s", targetBranch)
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

func ForcePush(remote string, branch string) string {
	return fmt.Sprintf("git push --force-with-lease -u %s %s", remote, branch)
}

func PushSetUpstream(remote string, branch string) string {
	return fmt.Sprintf("git push -u %s %s", remote, branch)
}

func BranchSetUpstream(remote string, branch string) string {
	return fmt.Sprintf("git branch -u %s/%s", remote, branch)
}

func Fetch(remote string) string {
	return fmt.Sprintf("git fetch %s", remote)
}

func MergeBase(branch string, parentBranch string) string {
	return fmt.Sprintf("git merge-base %s %s", parentBranch, branch)
}

func GetCommitHash(branch string) string {
	return fmt.Sprintf("git rev-parse %s", branch)
}

func ResetRemote(remote string, branch string) string {
	return fmt.Sprintf("git reset --hard %s/%s", remote, branch)
}

func GitHubPrCreate(parentBranch string) string {
	return fmt.Sprintf("gh pr create --base %s --fill", parentBranch)
}

func GitHubPrEditBody(number int, body string) string {
	return fmt.Sprintf("gh pr edit %d --body-file - <<EOF\n%s\nEOF", number, body)
}

func GitHubPrState() string {
	return "gh pr status --json state --jq '.currentBranch.state'"
}

func GitHubPrNumber() string {
	return "gh pr status --json number --jq '.currentBranch.number'"
}
