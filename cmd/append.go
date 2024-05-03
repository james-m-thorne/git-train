package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "Create a new branch from the current one, and store the parent in config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newBranch := args[0]
		currentBranch := command.GetOutputFatal(git.GetCurrentBranch())
		RunFatal(git.ConfigSetParent(newBranch, currentBranch))
		RunFatal(git.CheckoutNewBranch(newBranch))
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
}
