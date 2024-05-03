/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
	RunE: func(cmd *cobra.Command, args []string) error {
		newBranch := args[0]
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if err != nil {
			return err
		}
		err = Run(git.ConfigSetParent(newBranch, currentBranch))
		if err != nil {
			return err
		}
		return Run(git.CheckoutNewBranch(newBranch))
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
}
