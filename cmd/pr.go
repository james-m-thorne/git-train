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
		currentBranch := command.GetOutputFatal(git.GetCurrentBranch())
		if currentBranch == "" {
			command.PrintFatalError("current branch not found")
		}
		branchStack := []string{currentBranch}
		createParents, _ := cmd.Flags().GetBool("create-parents")
		if createParents {
			branchStack = git.GetBranchParentStack(currentBranch, false)
		}

		for _, branch := range branchStack {
			state, _ := command.GetOutput(git.GitHubPrState())
			if state == "" {
				RunFatal(git.GitHubPrCreate(branch))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().BoolP("create-parents", "c", false, "Create/update the PR's of the parent branches")
}
