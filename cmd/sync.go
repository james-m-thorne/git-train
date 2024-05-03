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

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all of the parent branches with upstream and to your current one",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			return fmt.Errorf("current branch not found")
		}

		includeMaster, _ := cmd.Flags().GetBool("include-master")
		branchStack := git.GetBranchParentStack(currentBranch, includeMaster)
		if len(branchStack) <= 1 {
			return fmt.Errorf("no parent branches found")
		}

		noUpdate, _ := cmd.Flags().GetBool("no-update")
		if !noUpdate {
			if err = Run(git.Checkout(branchStack[len(branchStack)-1])); err != nil {
				return fmt.Errorf("checkout failed: %s", err)
			}
			if err = Run(git.Pull()); err != nil {
				return fmt.Errorf("pull failed: %s", err)
			}
		}
		for i := len(branchStack) - 1; i >= 1; i-- {
			if err = Run(git.Checkout(branchStack[i-1])); err != nil {
				return fmt.Errorf("checkout failed: %s", err)
			}
			if !noUpdate {
				if err = Run(git.Pull()); err != nil {
					return fmt.Errorf("pull failed: %s", err)
				}
			}
			if err = Run(git.Rebase(branchStack[i])); err != nil {
				return fmt.Errorf("rebase failed: %s", err)
			}
			if !noUpdate {
				if err = Run(git.Push()); err != nil {
					return fmt.Errorf("push failed: %s", err)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP("include-master", "i", false, "Sync all the parent branches and include the master branch")
	syncCmd.Flags().BoolP("no-update", "n", false, "Do not pull/push the latest changes to remote vcs")
}
