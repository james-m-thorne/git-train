/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append",
	Short: "Create a new branch and append it to the current one",
	Long:  `Create a new branch and append it to the current one and store the parent in config`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		newBranch := args[0]
		currentBranch, err := command.GetOutput(git.GetCurrentBranch())
		if err != nil {
			return err
		}
		err = command.Exec(git.CheckoutNewBranch(newBranch))
		if err != nil {
			return err
		}
		return command.Exec(git.ConfigSetParent(currentBranch, newBranch))
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
