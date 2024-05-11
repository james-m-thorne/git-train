package cmd

import (
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

		branchStack := git.GetBranchStack(currentBranch, true)
		if len(branchStack) <= 1 {
			command.PrintFatalError("no parent branches found")
		}

		branchesToCreate := 1
		createAll, _ := cmd.Flags().GetBool("all")
		if createAll {
			branchesToCreate = len(branchStack)
		}

		prs := git.GetBranchStackPullRequests(branchStack)
		prs = git.UpdatePullRequestBodies(branchStack, prs)
		for i := 0; i < branchesToCreate; i++ {
			branch := branchStack[i]
			state, _ := command.GetOutput(git.GitHubPrState())
			if state == "" {
				RunFatal(git.Checkout(branch))
				RunFatal(git.PushSetUpstream(remote))
				parentBranch := branchStack[i+1]
				RunFatal(git.GitHubPrCreate(parentBranch))
			}
			RunFatal(git.GitHubPrEditBody(prs[branch].Number, prs[branch].Body))
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().BoolP("all", "a", false, "Create/update the PR's of the parent and children branches")
}
