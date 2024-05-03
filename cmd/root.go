/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-train",
	Short: "A way to manage stacked branches in git and GitHub",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var IsDryRun bool

func Run(shell string) error {
	return command.Exec(shell, IsDryRun)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&IsDryRun, "dry-run", "D", false, "defaultValue")
}
