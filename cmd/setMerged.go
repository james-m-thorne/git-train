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

// setMergedCmd represents the merge command
var setMergedCmd = &cobra.Command{
	Use:   "set-merged",
	Short: "Remove a branch train and rebase all of the descendants",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			return fmt.Errorf("current branch not found")
		}

		childBranch, err := command.GetOutput(git.ConfigGetChild(currentBranch))
		if childBranch != "" || err == nil {
			return fmt.Errorf("must be a branch with no children")
		}

		mergedBranch := args[0]
		if skipMergeCheck, _ := cmd.Flags().GetBool("skip-merge-check"); !skipMergeCheck {
			if err = Run(git.Checkout(mergedBranch)); err != nil {
				return fmt.Errorf("checkout failed: %s", err)
			}
			state, err := command.GetOutput(git.GitHubPrState())
			if state != "MERGED" || err != nil {
				return fmt.Errorf("parent branch is not merged on GitHub, state=%s", state)
			}
		}

		includeMaster, _ := cmd.Flags().GetBool("include-master")
		branchStack := git.GetBranchParentStack(currentBranch, includeMaster)

		hasPassedMergedBranch := false
		for i := len(branchStack) - 1; i >= 1; i-- {
			parentBranch := branchStack[i]
			currentBranch = branchStack[i-1]
			if parentBranch == mergedBranch {
				continue
			}
			if currentBranch == mergedBranch {
				if i-2 < 0 {
					return fmt.Errorf("merged branch has no parent: %s", mergedBranch)
				}
				currentBranch = branchStack[i-2]
				err = Run(git.ConfigSetParent(currentBranch, parentBranch))
				if err != nil {
					return fmt.Errorf("failed to set new parent branch for %s", currentBranch)
				}
				hasPassedMergedBranch = true
			}

			if err = Run(git.Checkout(currentBranch)); err != nil {
				return fmt.Errorf("checkout failed: %s", err)
			}
			if hasPassedMergedBranch {
				err = Run(git.RebaseOntoTarget(parentBranch, mergedBranch, currentBranch))
			} else {
				err = Run(git.Rebase(parentBranch))
			}
			if err != nil {
				return fmt.Errorf("rebase error: %s", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setMergedCmd)
	setMergedCmd.Flags().BoolP("skip-merge-check", "S", false, "Skip the merge check for the parent branch")
	setMergedCmd.Flags().BoolP("include-master", "i", true, "Sync all the parent branches and include the master branch")
}
