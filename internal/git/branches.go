package git

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/xlab/treeprint"
	"regexp"
	"strings"
)

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

func PrintBranchChildTree() {
	masterBranch, _ := command.GetOutput(ConfigGetMaster())

	tree := treeprint.NewWithRoot(masterBranch)
	AddChildBranches(tree, masterBranch)
	fmt.Println(tree.String())
}

func AddChildBranches(tree treeprint.Tree, branch string) {
	pattern := regexp.MustCompile(`git-train\.(.*?)\.parent`)

	childTree := tree.AddBranch(branch)
	childrenString, _ := command.GetOutput(ConfigGetChild(branch))
	children := strings.Split(childrenString, "\n")
	for _, child := range children {
		matches := pattern.FindStringSubmatch(child)
		if matches != nil && len(matches) > 1 {
			AddChildBranches(childTree, matches[1])
		}
	}
}
