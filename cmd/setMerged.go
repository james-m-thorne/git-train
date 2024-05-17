package cmd

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

// setMergedCmd represents the merge command
var setMergedCmd = &cobra.Command{
	Use:   "set-merged",
	Short: "Remove a branch train and rebase all of the descendants",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		_ = syncCmd.Flags().Set("validate", "true")
		_ = syncCmd.Flags().Set("no-update", "true")
		_ = syncCmd.Flags().Set("fetch", "true")
		_ = syncCmd.Flags().Set("merge", "true")
		_ = syncCmd.Flags().Set("push", "true")
		syncCmd.Run(syncCmd, []string{})

		remote := command.GetOutputFatal(git.ConfigGetRemote())
		currentBranch := git.GetCurrentBranch()

		mergedBranch := args[0]
		if skipMergeCheck, _ := cmd.Flags().GetBool("skip-merge-check"); !skipMergeCheck {
			RunFatal(git.Checkout(mergedBranch))
			state := command.GetOutputFatal(git.GitHubPrState())
			if state != "MERGED" {
				command.PrintFatalError("parent branch is not merged on GitHub, state=%s", state)
			}
		}

		excludeMaster, _ := cmd.Flags().GetBool("exclude-master")
		branchStack := git.GetBranchStack(currentBranch, excludeMaster)
		skipValidation, _ := cmd.Flags().GetBool("skip-validation")
		if !skipValidation {
			if !slices.Contains(branchStack, mergedBranch) {
				command.PrintFatalError("invalid branch: %s\nmust be one of %v", mergedBranch, branchStack)
			}
			git.ValidateBranchStack(branchStack, []string{mergedBranch})
		}

		completedBranchesStr := command.GetOutputFatal(git.ConfigGetMergedCompletedBranches(mergedBranch))
		var completedBranches []string
		if len(completedBranchesStr) > 0 {
			completedBranches = strings.Split(completedBranchesStr, ",")
		}

		updateParentCommand := ""
		for i := 1; i < len(branchStack); i++ {
			parentBranch := branchStack[i-1]
			currentBranch = branchStack[i]
			if currentBranch == mergedBranch || slices.Contains(completedBranches, currentBranch) {
				continue
			}
			if parentBranch == mergedBranch {
				if i-2 < 0 {
					command.PrintFatalError("%s does not have a valid parent branch", mergedBranch)
				}

				parentBranch = branchStack[i-2]
				skipUpdateParent, _ := cmd.Flags().GetBool("skip-update-parent")
				if !skipUpdateParent {
					updateParentCommand = git.ConfigSetParent(currentBranch, parentBranch)
				}
			}

			RunFatal(git.Checkout(currentBranch))
			RunFatal(git.RebaseOntoTarget(parentBranch, fmt.Sprintf("%s/%s", remote, parentBranch), currentBranch))

			completedBranches = append(completedBranches, currentBranch)
			Run(git.ConfigSetMergedCompletedBranches(mergedBranch, completedBranches))
		}

		if updateParentCommand != "" {
			RunFatal(updateParentCommand)
		}

		Run(git.ConfigDeleteMergedCompletedBranches(mergedBranch))
	},
}

func init() {
	rootCmd.AddCommand(setMergedCmd)
	setMergedCmd.Flags().BoolP("skip-validation", "v", false, "Skip the validation of branch stack")
	setMergedCmd.Flags().BoolP("skip-update-parent", "p", false, "Skip the ref update for the parent branch")
	setMergedCmd.Flags().BoolP("skip-merge-check", "m", false, "Skip the merge check for the parent branch")
	setMergedCmd.Flags().BoolP("exclude-master", "e", false, "Sync all the parent branches but exclude the master branch")
}
