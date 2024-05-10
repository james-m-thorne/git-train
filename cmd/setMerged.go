package cmd

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"github.com/spf13/cobra"
	"slices"
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
		currentBranch := command.GetOutputFatal(git.GetCurrentBranch())
		if currentBranch == "" {
			command.PrintFatalError("current branch not found")
		}

		children := git.GetBranchChildren(currentBranch)
		if len(children) > 0 {
			command.PrintFatalError("must be a branch with no children. try again after `git train last`")
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
		skipValidation, _ := cmd.Flags().GetBool("skip-validation")
		if !skipValidation {
			if !slices.Contains(branchStack, mergedBranch) {
				command.PrintFatalError("invalid branch: %s\nmust be one of %v", mergedBranch, branchStack)
			}
			git.ValidateBranchStack(branchStack, []string{mergedBranch})
		}

		updateParentCommand := ""
		for i := len(branchStack) - 1; i >= 1; i-- {
			parentBranch := branchStack[i]
			currentBranch = branchStack[i-1]
			if currentBranch == mergedBranch {
				continue
			}
			if parentBranch == mergedBranch {
				if i+1 >= len(branchStack) {
					command.PrintFatalError("%s does not have a valid parent branch", mergedBranch)
				}

				parentBranch = branchStack[i+1]
				skipUpdateParent, _ := cmd.Flags().GetBool("skip-update-parent")
				if !skipUpdateParent {
					updateParentCommand = git.ConfigSetParent(currentBranch, parentBranch)
				}
			}

			RunFatal(git.Checkout(currentBranch))
			RunFatal(git.RebaseOntoTarget(parentBranch, fmt.Sprintf("%s/%s", remote, parentBranch), currentBranch))
		}

		if updateParentCommand != "" {
			RunFatal(updateParentCommand)
		}
	},
}

func init() {
	rootCmd.AddCommand(setMergedCmd)
	setMergedCmd.Flags().BoolP("skip-validation", "v", false, "Skip the validation of branch stack")
	setMergedCmd.Flags().BoolP("skip-update-parent", "p", false, "Skip the ref update for the parent branch")
	setMergedCmd.Flags().BoolP("skip-merge-check", "m", false, "Skip the merge check for the parent branch")
	setMergedCmd.Flags().BoolP("exclude-master", "e", false, "Sync all the parent branches but exclude the master branch")
}
