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
		var branchStack []string
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			return fmt.Errorf("current branch not found")
		}

		masterBranch := ""
		if includeMaster, _ := cmd.Flags().GetBool("include-master"); !includeMaster {
			masterBranch, _ = command.GetOutput(git.ConfigGetMaster())
		}
		for currentBranch != masterBranch {
			branchStack = append(branchStack, currentBranch)
			currentBranch, err = command.GetOutput(git.ConfigGetParent(currentBranch))
			if err != nil {
				break
			}
		}

		for i := len(branchStack) - 1; i >= 1; i-- {
			if err = command.Exec(git.Checkout(branchStack[i-1])); err != nil {
				return fmt.Errorf("checkout failed: %s", err)
			}
			if err = command.Exec(git.Rebase(branchStack[i])); err != nil {
				return fmt.Errorf("rebase failed: %s", err)
			}
			err = command.Exec(git.GitHubPrState())
			if err == nil {
				if err = command.Exec(git.Push()); err != nil {
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
