package cmd

import (
	"fmt"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

func init() {
	pushFlag := false
	pullFlag := false

	cmd := &cobra.Command{
		Use:   "branch:current",
		Short: "Work with the current branch or return the current branch name",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			branchName, _ := helpers.GetCurrentBranchName()

			if pushFlag {
				helpers.RunCommandOnStdout("git", "push", "origin", branchName)
				return
			}

			if pullFlag {
				helpers.RunCommandOnStdout("git", "pull", "origin", branchName, "--rebase")
				return
			}

			fmt.Println(branchName)
		},
	}

	cmd.Flags().BoolVarP(&pushFlag, "push", "p", false, "push the current branch to origin")
	cmd.Flags().BoolVarP(&pullFlag, "pull", "u", false, "pull the current branch from origin and rebase")

	rootCmd.AddCommand(cmd)
}
