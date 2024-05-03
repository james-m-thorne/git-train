/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"log"

	"github.com/spf13/cobra"
)

// setParentCmd represents the setParent command
var setParentCmd = &cobra.Command{
	Use:   "set-parent",
	Short: "Sets the new parent branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newParentBranch := args[0]
		if err := Run(git.CheckBranchExists(newParentBranch)); err != nil {
			log.Fatalf("branch does not exists %s", newParentBranch)
		}

		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if err != nil {
			log.Fatalf("unable to get current branch")
		}
		masterBranch, _ := command.GetOutput(git.ConfigGetMaster())
		if currentBranch == masterBranch {
			return
		}

		if rebase, _ := cmd.Flags().GetBool("rebase"); rebase {
			oldParentBranch, err := command.GetOutput(git.ConfigGetParent(currentBranch))
			if err != nil {
				log.Fatalf("unable to get current branch")
			}
			err = Run(git.RebaseOntoTarget(newParentBranch, oldParentBranch, currentBranch))
			if err != nil {
				log.Fatalf("rebase failed, fix it and rerun this command")
			}
		}
		err = Run(git.ConfigSetParent(currentBranch, newParentBranch))
		if err != nil {
			log.Fatalf("failed to set parent: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setParentCmd)
	setParentCmd.Flags().BoolP("rebase", "r", false, "Rebase on the new parent")
}
