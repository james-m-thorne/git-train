/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config values for git-train",
	Long:  `Get all config values for git-train`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.Exec(git.ConfigGetAll())
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config values for git-train",
	Long: `Set values for:
master_branch: Name of the master branch`,
	Args: cobra.ExactArgs(2),
}

var configSetMasterBranchCmd = &cobra.Command{
	Use:   "master_branch",
	Short: "Set the master branch name",
	Long:  `Set the default name of the master branch`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return command.Exec(git.ConfigSetMaster(args[0]))
	},
}

func init() {
	configSetCmd.AddCommand(configSetMasterBranchCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
