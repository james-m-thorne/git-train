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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if currentBranch == "" || err != nil {
			return fmt.Errorf("current branch not found")
		}
		parentBranch, err := command.GetOutput(git.ConfigGetParent(currentBranch))
		if parentBranch == "" || err != nil {
			return fmt.Errorf("no parent branch found for %s", currentBranch)
		}
		parentsParentBranch, err := command.GetOutput(git.ConfigGetParent(parentBranch))
		if parentsParentBranch == "" || err != nil {
			return fmt.Errorf("no parent branch found for %s", parentBranch)
		}
		err = command.Exec(git.RebaseOntoParent(parentsParentBranch, parentBranch, currentBranch))
		if err != nil {
			return fmt.Errorf("rebase error: %s", err)
		}
		_ = command.Exec(git.ConfigDeleteParent(parentBranch))
		return command.Exec(git.Delete(parentBranch))
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
	// rebaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
