/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"

	"github.com/spf13/cobra"
)

// setParentCmd represents the setParent command
var setParentCmd = &cobra.Command{
	Use:   "set-parent",
	Short: "Sets the new parent branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newParentBranch := args[0]
		RunFatal(git.CheckBranchExists(newParentBranch))

		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if err != nil {
			command.PrintFatalError("unable to get current branch")
		}
		masterBranch, _ := command.GetOutput(git.ConfigGetMaster())
		if currentBranch == masterBranch {
			return
		}

		if rebase, _ := cmd.Flags().GetBool("rebase"); rebase {
			oldParentBranch, err := command.GetOutput(git.ConfigGetParent(currentBranch))
			if err != nil {
				command.PrintFatalError("unable to get current branch")
			}
			RunFatal(git.RebaseOntoTarget(newParentBranch, oldParentBranch, currentBranch))
			if err != nil {
				command.PrintFatalError("rebase failed, fix it and rerun this command")
			}
		}
		RunFatal(git.ConfigSetParent(currentBranch, newParentBranch))
	},
}

func init() {
	rootCmd.AddCommand(setParentCmd)
	setParentCmd.Flags().BoolP("rebase", "r", false, "Rebase on the new parent")
}
