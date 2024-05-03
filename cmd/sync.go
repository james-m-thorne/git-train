/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"log"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all of the parent branches with upstream and to your current one",
	Run: func(cmd *cobra.Command, args []string) {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			log.Fatalf("current branch not found")
		}

		includeMaster, _ := cmd.Flags().GetBool("include-master")
		branchStack := git.GetBranchParentStack(currentBranch, includeMaster)
		if len(branchStack) <= 1 {
			log.Fatalf("no parent branches found")
		}

		noUpdate, _ := cmd.Flags().GetBool("no-update")
		if !noUpdate {
			if err = Run(git.Checkout(branchStack[len(branchStack)-1])); err != nil {
				log.Fatalf("checkout failed: %s", err)
			}
			if err = Run(git.Pull()); err != nil {
				log.Fatalf("pull failed: %s", err)
			}
		}
		for i := len(branchStack) - 1; i >= 1; i-- {
			if err = Run(git.Checkout(branchStack[i-1])); err != nil {
				log.Fatalf("checkout failed: %s", err)
			}
			if !noUpdate {
				if err = Run(git.Pull()); err != nil {
					log.Fatalf("pull failed: %s", err)
				}
			}
			if err = Run(git.Rebase(branchStack[i])); err != nil {
				log.Fatalf("rebase failed: %s", err)
			}
			if !noUpdate {
				if err = Run(git.Push()); err != nil {
					log.Fatalf("push failed: %s", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP("include-master", "i", false, "Sync all the parent branches and include the master branch")
	syncCmd.Flags().BoolP("no-update", "n", false, "Do not pull/push the latest changes to remote vcs")
}
