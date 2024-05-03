/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"

	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create a GitHub PR on the current branch and set the base as the parent",
	Run: func(cmd *cobra.Command, args []string) {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			command.PrintFatalError("current branch not found")
		}
		parentBranch, err := command.GetOutput(git.ConfigGetParent(currentBranch))
		if parentBranch == "" || err != nil {
			command.PrintFatalError("no parent branch found for %s", currentBranch)
		}
		RunFatal(git.GitHubPrCreate(parentBranch))
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
}
