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
	Short: "Sync all of the parent branches with upstream and to your current one",
	Run: func(cmd *cobra.Command, args []string) {
		remote := command.GetOutputFatal(git.ConfigGetRemote())
		currentBranch := command.GetOutputFatal(git.GetCurrentBranch())
		if currentBranch == "" {
			command.PrintFatalError("current branch not found")
		}

		excludeMaster, _ := cmd.Flags().GetBool("exclude-master")
		branchStack := git.GetBranchParentStack(currentBranch, excludeMaster)
		if len(branchStack) <= 1 {
			command.PrintFatalError("no parent branches found")
		}

		shouldFetch, _ := cmd.Flags().GetBool("fetch")
		shouldMerge, _ := cmd.Flags().GetBool("merge")
		shouldPush, _ := cmd.Flags().GetBool("push")
		shouldValidate, _ := cmd.Flags().GetBool("validate")
		noUpdate, _ := cmd.Flags().GetBool("no-update")

		if shouldFetch {
			RunFatal(git.Fetch(remote))
		}
		for i := len(branchStack) - 1; i >= 1; i-- {
			currentBranch := branchStack[i-1]
			RunFatal(git.Checkout(currentBranch))
			if shouldMerge {
				RunFatal(git.Merge(fmt.Sprintf("%s/%s", remote, currentBranch)))
			}
			if !noUpdate {
				parentBranch := branchStack[i]
				RunFatal(git.RebaseOntoTarget(parentBranch, fmt.Sprintf("%s/%s", remote, parentBranch), currentBranch))
			}
			if shouldPush {
				RunFatal(git.ForcePush(remote, currentBranch))
			}
			if shouldValidate {
				git.CheckInSyncWithRemoteBranch(remote, currentBranch)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolP("exclude-master", "e", false, "Sync all the parent branches and exclude the master branch")
	syncCmd.Flags().BoolP("validate", "v", false, "Validate the branches are in sync with remote")
	syncCmd.Flags().BoolP("fetch", "f", true, "Fetch the latest changes from remote vcs")
	syncCmd.Flags().BoolP("merge", "m", true, "Merge the changes from remote vcs")
	syncCmd.Flags().BoolP("push", "p", false, "Push the latest changes to remote vcs")
	syncCmd.Flags().BoolP("no-update", "n", false, "Do not rebase with the parent branch")
}
