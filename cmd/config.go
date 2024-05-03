package cmd

import (
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage config values for git-train",
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config values for git-train",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		RunFatal(git.ConfigGetAll())
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config values for git-train",
	Args:  cobra.ExactArgs(2),
}

var configSetMasterBranchCmd = &cobra.Command{
	Use:   "master_branch",
	Short: "Set the master branch name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		RunFatal(git.ConfigSetMaster(args[0]))
	},
}

func init() {
	configSetCmd.AddCommand(configSetMasterBranchCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
