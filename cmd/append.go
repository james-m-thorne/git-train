/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
	"log"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "Create a new branch from the current one, and store the parent in config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newBranch := args[0]
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if err != nil {
			log.Fatalf("failed to set parent %s", currentBranch)
		}
		err = Run(git.ConfigSetParent(newBranch, currentBranch))
		if err != nil {
			log.Fatalf("failed to set parent %s", currentBranch)
		}
		err = Run(git.CheckoutNewBranch(newBranch))
		if err != nil {
			log.Fatalf("checkout failed for %s", newBranch)
		}
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
}
