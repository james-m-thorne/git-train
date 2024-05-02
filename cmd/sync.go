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
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	syncCmd.Flags().BoolP("include-master", "i", false, "Sync all the parent branches and include the master branch")
}
