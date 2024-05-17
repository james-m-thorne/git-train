package git

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/xlab/treeprint"
	"regexp"
	"slices"
	"strings"
)

func GetCurrentBranch() string {
	branch := command.GetOutputFatal(ShowCurrentBranch())
	if branch == "" {
		command.PrintFatalError("current branch name is empty")
	}
	return branch
}

func GetBranchChildren(branch string) []string {
	pattern := regexp.MustCompile(`git-train\.(.*?)\.parent`)
	childrenString, _ := command.GetOutput(ConfigGetChild(branch))
	childrenConfigValues := strings.Split(childrenString, "\n")

	var children []string
	for _, child := range childrenConfigValues {
		matches := pattern.FindStringSubmatch(child)
		if matches != nil && len(matches) > 1 {
			childBranch := matches[1]
			command.GetOutputFatal(CheckBranchExists(childBranch))
			children = append(children, childBranch)
		}
	}
	return children
}

// GetBranchChildStack return the current branch and all its children
func GetBranchChildStack(currentBranch string) []string {
	branches := []string{currentBranch}
	children := GetBranchChildren(currentBranch)
	for _, branch := range children {
		branches = append(branches, GetBranchChildStack(branch)...)
	}
	return branches
}

// GetBranchParentStack return the current branch and all its parents
func GetBranchParentStack(currentBranch string, excludeMaster bool) []string {
	var branchStack []string
	masterBranch := ""
	if excludeMaster {
		masterBranch, _ = command.GetOutput(ConfigGetMaster())
	}
	for currentBranch != masterBranch {
		command.GetOutputFatal(CheckBranchExists(currentBranch))

		branchStack = append(branchStack, currentBranch)
		parentBranch, err := command.GetOutput(ConfigGetParent(currentBranch))
		if err != nil {
			break
		}
		currentBranch = parentBranch
	}
	return branchStack
}

func GetBranchStack(currentBranch string, excludeMaster bool) []string {
	parentStack := GetBranchParentStack(currentBranch, excludeMaster)
	slices.Reverse(parentStack)
	childStack := GetBranchChildStack(currentBranch)
	return append(parentStack, childStack[1:]...)
}

func ValidateBranchStack(branchStack []string, skipValidationForBranches []string) {
	skipBranchesSet := make(map[string]bool) // Create a map to represent the set
	for _, item := range skipValidationForBranches {
		skipBranchesSet[item] = true // Add each item to the set
	}

	for i := 1; i < len(branchStack); i++ {
		parentBranch := branchStack[i-1]
		currentBranch := branchStack[i]
		if _, ok := skipBranchesSet[currentBranch]; !ok {
			mergeHash := command.GetOutputFatal(MergeBase(currentBranch, parentBranch))
			parentHeadHash := command.GetOutputFatal(GetCommitHash(parentBranch))
			if mergeHash != parentHeadHash {
				command.PrintFatalError("non-linear branches for parent branch %s: %s and current branch %s : %s. try sync and rebase the branches", parentBranch, parentHeadHash, currentBranch, mergeHash)
			}
		}
	}
}

func CheckInSyncWithRemoteBranch(remote string, branch string) {
	currentHeadHash := command.GetOutputFatal(GetCommitHash(branch))
	remoteHeadHash := command.GetOutputFatal(GetCommitHash(fmt.Sprintf("%s/%s", remote, branch)))
	if currentHeadHash != remoteHeadHash {
		command.PrintFatalError("%s is not in sync with remote", branch)
	}
}

func AddChildBranches(tree treeprint.Tree, branch string) {
	childTree := tree.AddBranch(branch)
	children := GetBranchChildren(branch)
	for _, child := range children {
		AddChildBranches(childTree, child)
	}
}

func GetReadableCommitHash(remote string, branch string) string {
	remoteBranch := fmt.Sprintf("%s/%s", remote, branch)
	branchHash := command.GetOutputFatal(GetCommitHash(branch))
	remoteBranchHash := command.GetOutputFatal(GetCommitHash(remoteBranch))
	if branchHash == remoteBranchHash {
		return remoteBranch
	}
	return branchHash
}
