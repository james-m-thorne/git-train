/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
	"strings"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all of the parent branches with upstream and to your current one",
	Run: func(cmd *cobra.Command, args []string) {
		currentBranch := command.GetOutputFatal(git.GetCurrentBranch())
		if currentBranch == "" {
			command.PrintFatalError("current branch not found")
		}

		includeMaster, _ := cmd.Flags().GetBool("include-master")
		branchStack := git.GetBranchParentStack(currentBranch, includeMaster)
		if len(branchStack) <= 1 {
			command.PrintFatalError("no parent branches found")
		}

		strategy, _ := cmd.Flags().GetString("strategy")
		noUpdate, _ := cmd.Flags().GetBool("no-update")
		if !noUpdate {
			RunFatal(git.Checkout(branchStack[len(branchStack)-1]))
			RunFatal(git.Pull())
		}
		for i := len(branchStack) - 1; i >= 1; i-- {
			RunFatal(git.Checkout(branchStack[i-1]))
			if !noUpdate {
				RunFatal(git.Pull())
			}
			if strings.ToLower(strategy) == "merge" {
				RunFatal(git.Merge(branchStack[i]))
			} else {
				RunFatal(git.Rebase(branchStack[i]))
			}
			if !noUpdate {
				RunFatal(git.Push())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP("include-master", "i", false, "Sync all the parent branches and include the master branch")
	syncCmd.Flags().BoolP("no-update", "n", false, "Do not pull/push the latest changes to remote vcs")
	syncCmd.Flags().StringP("strategy", "s", "rebase", "Sync strategy for branches. Either merge or rebase")
}
