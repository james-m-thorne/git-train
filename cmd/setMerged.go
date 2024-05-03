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

		includeMaster, _ := cmd.Flags().GetBool("include-master")
		branchStack := git.GetBranchParentStack(currentBranch, includeMaster)

		hasPassedMergedBranch := false
		for i := len(branchStack) - 1; i >= 1; i-- {
			parentBranch := branchStack[i]
			currentBranch = branchStack[i-1]
			if parentBranch == mergedBranch {
				continue
			}
			if currentBranch == mergedBranch {
				if i-2 < 0 {
					command.PrintFatalError("merged branch has no parent: %s", mergedBranch)
				}
				currentBranch = branchStack[i-2]
				RunFatal(git.ConfigSetParent(currentBranch, parentBranch))
				skipPull, _ := cmd.Flags().GetBool("skip-pull")
				if !skipPull {
					RunFatal(git.Checkout(parentBranch))
					RunFatal(git.Pull())
				}
				hasPassedMergedBranch = true
			}

			RunFatal(git.Checkout(currentBranch))
			if hasPassedMergedBranch {
				RunFatal(git.RebaseOntoTarget(parentBranch, mergedBranch, currentBranch))
			} else {
				RunFatal(git.Rebase(parentBranch))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(setMergedCmd)
	setMergedCmd.Flags().BoolP("skip-pull", "p", false, "Skip the pull for the parent branch")
	setMergedCmd.Flags().BoolP("skip-merge-check", "S", false, "Skip the merge check for the parent branch")
	setMergedCmd.Flags().BoolP("include-master", "i", true, "Sync all the parent branches and include the master branch")
}
