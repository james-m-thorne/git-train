package cmd

import (
	"github.com/james-m-thorne/git-train/internal/command"
	"github.com/james-m-thorne/git-train/internal/git"
	"strconv"

	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Create a GitHub PR on the current branch and set the base as the parent",
	Run: func(cmd *cobra.Command, args []string) {
		remote := command.GetOutputFatal(git.ConfigGetRemote())
		currentBranch := git.GetCurrentBranch()
		masterBranch := command.GetOutputFatal(git.ConfigGetMaster())
		if currentBranch == masterBranch {
			command.PrintFatalError("cannot run pr on %s branch", masterBranch)
		}

		var branchStack []string
		createAll, _ := cmd.Flags().GetBool("all")
		if createAll {
			branchStack = git.GetBranchStack(currentBranch, false)
			if len(branchStack) <= 1 {
				command.PrintFatalError("no parent branches found")
			}
		} else {
			parentBranch := command.GetOutputFatal(git.ConfigGetParent(currentBranch))
			branchStack = []string{parentBranch, currentBranch}
		}

		prs := git.GetBranchStackPullRequests(branchStack)
		prs = git.UpdatePullRequestBodies(branchStack[1:], prs)
		for i := 1; i < len(branchStack); i++ {
			parentBranch := branchStack[i-1]
			branch := branchStack[i]
			RunFatal(git.Checkout(branch))
			state, _ := command.GetOutput(git.GitHubPrState())
			if state == "" {
				RunFatal(git.PushSetUpstream(remote, currentBranch))
				RunFatal(git.GitHubPrCreate(parentBranch))
				prNumberStr := command.GetOutputFatal(git.GitHubPrNumber())
				prNumber, err := strconv.Atoi(prNumberStr)
				if err != nil {
					command.PrintFatalError("failed to parse PR number: %s", prNumberStr)
				}
				RunFatal(git.GitHubPrEditBody(prNumber, prs[branch].Body))
			} else if state == "OPEN" {
				RunFatal(git.GitHubPrEditBody(prs[branch].Number, prs[branch].Body))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.Flags().BoolP("all", "a", false, "Create/update the PR's of the parent and children branches")
}
