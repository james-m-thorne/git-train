/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a tree like view of your branches stored in config",
	Run: func(cmd *cobra.Command, args []string) {
		masterBranch, _ := command.GetOutput(git.ConfigGetMaster())

		tree := treeprint.New()
		git.AddChildBranches(tree, masterBranch)
		fmt.Println(tree.String())
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
