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

// rebaseCmd represents the merge command
var rebaseCmd = &cobra.Command{
	Use:   "rebase",
	Short: "A brief description of your command",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			return fmt.Errorf("current branch not found")
		}
		parentBranch, err := command.GetOutput(git.ConfigGetParent(currentBranch))
		if parentBranch == "" || err != nil {
			return fmt.Errorf("no parent branch found for %s", currentBranch)
		}
		if skipMergeCheck, _ := cmd.Flags().GetBool("skip-merge-check"); skipMergeCheck {
			state, err := command.GetOutput(git.GitHubPrState())
			if state != "MERGED" || err != nil {
				return fmt.Errorf("parent branch is not merged on GitHub, state=%s", state)
			}
		}
		parentsParentBranch, err := command.GetOutput(git.ConfigGetParent(parentBranch))
		if parentsParentBranch == "" || err != nil {
			return fmt.Errorf("no parent branch found for %s", parentBranch)
		}
		err = command.Exec(git.RebaseOntoTarget(parentsParentBranch, parentBranch, currentBranch))
		if err != nil {
			return fmt.Errorf("rebase error: %s", err)
		}
		return command.Exec(git.ConfigSetParent(currentBranch, parentBranch))
	},
}

func init() {
	rootCmd.AddCommand(rebaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rebaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rebaseCmd.Flags().BoolP("skip-merge-check", "S", false, "Skip the merge check for the parent branch")
}
