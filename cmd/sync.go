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
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			return fmt.Errorf("current branch not found")
		}

		includeMaster, _ := cmd.Flags().GetBool("include-master")
		branchStack := git.GetBranchStack(currentBranch, includeMaster)

		for i := len(branchStack) - 1; i >= 1; i-- {
			if err = Run(git.Checkout(branchStack[i-1])); err != nil {
				return fmt.Errorf("checkout failed: %s", err)
			}
			if err = Run(git.Rebase(branchStack[i])); err != nil {
				return fmt.Errorf("rebase failed: %s", err)
			}
			_, err := command.GetOutput(git.GitHubPrState())
			if err == nil {
				if err = Run(git.Push()); err != nil {
					return fmt.Errorf("push failed: %s", err)
				}
			} else {
				fmt.Printf("no remote found: skipping push for %s\n", branchStack[i])
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP("include-master", "i", false, "Sync all the parent branches and include the master branch")
}
