package git

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/xlab/treeprint"
	"regexp"
	"strings"
)

func GetBranchChildren(branch string) []string {
	pattern := regexp.MustCompile(`git-train\.(.*?)\.parent`)
	childrenString, _ := command.GetOutput(ConfigGetChild(branch))
	childrenConfigValues := strings.Split(childrenString, "\n")

	var children []string
	for _, child := range childrenConfigValues {
		matches := pattern.FindStringSubmatch(child)
		if matches != nil && len(matches) > 1 {
			children = append(children, matches[1])
		}
	}
	return children
}

func GetBranchParentStack(currentBranch string, includeMaster bool) []string {
	var branchStack []string
	masterBranch := ""
	if !includeMaster {
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

func ValidateBranchStack(branchStack []string, skipValidationForBranches []string) {
	skipBranchesSet := make(map[string]bool) // Create a map to represent the set
	for _, item := range skipValidationForBranches {
		skipBranchesSet[item] = true // Add each item to the set
	}

	for i := len(branchStack) - 1; i >= 1; i-- {
		parentBranch := branchStack[i]
		currentBranch := branchStack[i-1]
		if _, ok := skipBranchesSet[currentBranch]; !ok {
			mergeHash := command.GetOutputFatal(MergeBase(currentBranch, parentBranch))
			parentHeadHash := command.GetOutputFatal(GetCommitHash(parentBranch))
			if mergeHash != parentHeadHash {
				command.PrintFatalError("non-linear branches for parent branch %s: %s and current branch %s : %s. try sync and rebase the branches", parentBranch, parentHeadHash, currentBranch, mergeHash)
			}
		}
	}
}

func AddChildBranches(tree treeprint.Tree, branch string) {
	childTree := tree.AddBranch(branch)
	children := GetBranchChildren(branch)
	for _, child := range children {
		AddChildBranches(childTree, child)
	}
}
