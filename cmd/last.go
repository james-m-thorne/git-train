package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"

	"github.com/spf13/cobra"
)

// lastCmd represents the last command
var lastCmd = &cobra.Command{
	Use:   "last",
	Short: "Checkout the last branch in the train",
	Run: func(cmd *cobra.Command, args []string) {
		currentBranch := command.GetOutputFatal(git.GetCurrentBranch())
		if currentBranch == "" {
			command.PrintFatalError("current branch not found")
		}

		children := git.GetAllChildBranches(currentBranch)
		RunFatal(git.Checkout(children[len(children)-1]))
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
}
