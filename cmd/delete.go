package cmd

import (
	"bufio"
	"fmt"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// deleteCmd represents the remove command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a branch and remove the stored parent config",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		deleteBranch := args[0]
		deleteAll, _ := cmd.Flags().GetBool("all")
		if !deleteAll {
			return nil
		}
		deleteBranches := git.GetBranchStack(deleteBranch, true)
		fmt.Println(fmt.Sprintf("Branches to delete: %s", strings.Join(deleteBranches, ", ")))
		fmt.Println("Are you sure you want to delete these branches? (y/n)")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" {
			return fmt.Errorf("stopping delete")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		deleteBranch := args[0]
		RunFatal(git.CheckBranchExists(deleteBranch))

		deleteBranches := []string{deleteBranch}
		deleteAll, _ := cmd.Flags().GetBool("all")
		if deleteAll {
			deleteBranches = git.GetBranchStack(deleteBranch, true)
		}

		for _, branch := range deleteBranches {
			Run(git.ConfigDeleteParent(branch))
			RunFatal(git.Delete(branch))
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("all", "a", false, "Delete all the branches in the stack")
}
