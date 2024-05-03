/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/git"

	"github.com/spf13/cobra"
)

// deleteCmd represents the remove command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a branch and remove the stored parent config",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
		_ = Run(git.ConfigDeleteParent(branch))
		return Run(git.Delete(branch))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
