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
		branchStack := git.GetBranchParentStack(currentBranch, true)

		branchesToCreate := 1
		createAllParents, _ := cmd.Flags().GetBool("create-parents")
		if createAllParents {
			branchesToCreate = len(branchStack) - 1
		}

		for i := 0; i < branchesToCreate; i++ {
			branch := branchStack[i]
			parentBranch := branchStack[i+1]
			RunFatal(git.Checkout(branch))
			RunFatal(git.PushSetUpstream())
			state, _ := command.GetOutput(git.GitHubPrState())
			if state == "" {
				RunFatal(git.GitHubPrCreate(parentBranch))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().BoolP("create-parents", "c", false, "Create/update the PR's of the parent branches")
}
