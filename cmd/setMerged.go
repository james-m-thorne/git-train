package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
)

// setMergedCmd represents the merge command
var setMergedCmd = &cobra.Command{
	Use:   "set-merged",
	Short: "Remove a branch train and rebase all of the descendants",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		currentBranch := command.GetOutputFatal(git.GetCurrentBranch())
		if currentBranch == "" {
			command.PrintFatalError("current branch not found")
		}

		children := git.GetBranchChildren(currentBranch)
		if len(children) > 0 {
			command.PrintFatalError("must be a branch with no children. try again after `git checkout %s`", children[len(children)-1])
		}

		mergedBranch := args[0]
		if skipMergeCheck, _ := cmd.Flags().GetBool("skip-merge-check"); !skipMergeCheck {
			RunFatal(git.Checkout(mergedBranch))
			state := command.GetOutputFatal(git.GitHubPrState())
			if state != "MERGED" {
				command.PrintFatalError("parent branch is not merged on GitHub, state=%s", state)
			}
		}

		excludeMaster, _ := cmd.Flags().GetBool("exclude-master")
		branchStack := git.GetBranchParentStack(currentBranch, excludeMaster)
		git.ValidateBranchStack(branchStack, []string{mergedBranch})

		isMergedBranch := false
		updateParentCommand := ""
		mergeBaseHash := ""
		for i := len(branchStack) - 1; i >= 2; i-- {
			grandParentBranch := branchStack[i]
			parentBranch := branchStack[i-1]
			currentBranch = branchStack[i-2]

			if mergeBaseHash == "" {
				mergeBaseHash = command.GetOutputFatal(git.GetCommitHash(parentBranch))
			}

			if parentBranch == mergedBranch {
				isMergedBranch = true
				skipUpdateParent, _ := cmd.Flags().GetBool("skip-update-parent")
				if !skipUpdateParent {
					updateParentCommand = git.ConfigSetParent(currentBranch, grandParentBranch)
				}
				skipPull, _ := cmd.Flags().GetBool("skip-pull")
				if !skipPull {
					RunFatal(git.Checkout(parentBranch))
					RunFatal(git.Pull())
				}
				mergeBaseHash = command.GetOutputFatal(git.GetCommitHash(parentBranch))
			}

			RunFatal(git.Checkout(currentBranch))
			beforeRebaseMergeBaseHash := command.GetOutputFatal(git.GetCommitHash(currentBranch))
			if isMergedBranch {
				RunFatal(git.RebaseOntoTarget(grandParentBranch, parentBranch, currentBranch))
				isMergedBranch = false
			} else {
				RunFatal(git.RebaseOntoTarget(parentBranch, mergeBaseHash, currentBranch))
			}
			mergeBaseHash = beforeRebaseMergeBaseHash
		}

		if updateParentCommand != "" {
			RunFatal(updateParentCommand)
		}
	},
}

func init() {
	rootCmd.AddCommand(setMergedCmd)
	setMergedCmd.Flags().BoolP("skip-pull", "l", false, "Skip the pull for the parent branch")
	setMergedCmd.Flags().BoolP("skip-update-parent", "p", false, "Skip the update for the parent branch")
	setMergedCmd.Flags().BoolP("skip-merge-check", "S", false, "Skip the merge check for the parent branch")
	setMergedCmd.Flags().BoolP("exclude-master", "e", false, "Sync all the parent branches but exclude the master branch")
}
