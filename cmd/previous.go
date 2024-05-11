package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// previousCmd represents the previous command
var previousCmd = &cobra.Command{
	Use:     "previous",
	Aliases: []string{"p", "prev"},
	Short:   "Checkout the previous branch in the train",
	Run: func(cmd *cobra.Command, args []string) {
		currentBranch := git.GetCurrentBranch()
		parent := command.GetOutputFatal(git.ConfigGetParent(currentBranch))
		RunFatal(git.Checkout(parent))
	},
}

func init() {
	rootCmd.AddCommand(previousCmd)
}
