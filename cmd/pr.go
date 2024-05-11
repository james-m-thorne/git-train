package cmd

import (
	"fmt"
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"

	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create a GitHub PR on the current branch and set the base as the parent",
	Run: func(cmd *cobra.Command, args []string) {
		remote := command.GetOutputFatal(git.ConfigGetRemote())
		currentBranch := git.GetCurrentBranch()

		branchStack := git.GetBranchParentStack(currentBranch, true)
		if len(branchStack) <= 1 {
			command.PrintFatalError("no parent branches found")
		}

		branchesToCreate := 1
		createAllParents, _ := cmd.Flags().GetBool("create-parents")
		if createAllParents {
			branchesToCreate = len(branchStack) - 1
		}

		for i := 0; i < branchesToCreate; i++ {
			branch := branchStack[i]
			parentBranch := branchStack[i+1]
			RunFatal(git.Checkout(branch))
			RunFatal(git.PushSetUpstream(remote))
			state, _ := command.GetOutput(git.GitHubPrState())
			if state == "" {
				RunFatal(git.GitHubPrCreate(parentBranch))
			} else {
				RunFatal(git.GitHubPrView())
			}
			body := command.GetOutputFatal(git.GitHubPrBody())
			fmt.Println(body)
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().BoolP("create-parents", "c", false, "Create/update the PR's of the parent branches")
}
