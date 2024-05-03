/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"

	"github.com/spf13/cobra"
)

// setParentCmd represents the setParent command
var setParentCmd = &cobra.Command{
	Use:   "set-parent",
	Short: "Sets the new parent branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		newParentBranch := args[0]
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if err != nil {
			return fmt.Errorf("unable to get current branch")
		}
		if rebase, _ := cmd.Flags().GetBool("rebase"); rebase {
			oldParentBranch, err := command.GetOutput(git.ConfigGetParent(currentBranch))
			if err != nil {
				return fmt.Errorf("unable to get current branch")
			}
			err = Run(git.RebaseOntoTarget(newParentBranch, oldParentBranch, currentBranch))
			if err != nil {
				return fmt.Errorf("rebase failed, fix it and rerun this command")
			}
		}
		return Run(git.ConfigSetParent(currentBranch, newParentBranch))
	},
}

func init() {
	rootCmd.AddCommand(setParentCmd)
	setParentCmd.Flags().BoolP("rebase", "r", false, "Rebase on the new parent")
}
