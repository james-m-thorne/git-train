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
		branchStack = append(branchStack, currentBranch)
		parentBranch, err := command.GetOutput(ConfigGetParent(currentBranch))
		if err != nil {
			break
		}
		currentBranch = parentBranch
	}
	return branchStack
}

func AddChildBranches(tree treeprint.Tree, branch string) {
	childTree := tree.AddBranch(branch)
	children := GetBranchChildren(branch)
	for _, child := range children {
		AddChildBranches(childTree, child)
	}
}
