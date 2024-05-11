package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all parent branches on remote branches",
	Run: func(cmd *cobra.Command, args []string) {
		remote := command.GetOutputFatal(git.ConfigGetRemote())
		currentBranch := git.GetCurrentBranch()
		branchStack := git.GetBranchParentStack(currentBranch, true)
		for i := len(branchStack) - 1; i >= 0; i-- {
			RunFatal(git.Checkout(branchStack[i]))
			RunFatal(git.ResetRemote(remote, branchStack[i]))
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
