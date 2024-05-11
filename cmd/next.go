package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:     "next",
	Aliases: []string{"n"},
	Short:   "Checkout the next branch in the train",
	Run: func(cmd *cobra.Command, args []string) {
		currentBranch := git.GetCurrentBranch()
		children := git.GetBranchChildren(currentBranch)
		if len(children) == 0 {
			command.PrintFatalError("%s has no child branch", currentBranch)
		} else if len(children) > 1 {
			command.PrintFatalError("%s has more than one child", currentBranch)
		}

		child := children[0]
		RunFatal(git.Checkout(child))
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
