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

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			return fmt.Errorf("current branch not found")
		}
		parentBranch, err := command.GetOutput(git.ConfigGetParent(currentBranch))
		if parentBranch == "" || err != nil {
			return fmt.Errorf("no parent branch found for %s", currentBranch)
		}
		return Run(git.GitHubPrCreate(parentBranch))
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
}
