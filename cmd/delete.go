/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// deleteCmd represents the remove command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a branch and remove the stored parent config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteBranch := args[0]
		RunFatal(git.CheckBranchExists(deleteBranch))

		deleteBranches := []string{deleteBranch}
		deleteChildren, _ := cmd.Flags().GetBool("children")
		if deleteChildren {
			deleteBranches = addChildBranches(deleteBranches, deleteBranch)
		}

		for _, branch := range deleteBranches {
			Run(git.ConfigDeleteParent(branch))
			RunFatal(git.Delete(branch))
		}
	},
}

func addChildBranches(branches []string, currentBranch string) []string {
	children := git.GetBranchChildren(currentBranch)
	branches = append(branches, children...)
	for _, branch := range children {
		branches = addChildBranches(branches, branch)
	}
	return branches
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("children", "c", false, "Delete all the children of this branch")
}
