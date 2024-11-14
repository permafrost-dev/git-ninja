package cmd

import (
	"fmt"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

func init() {
	flagRemoteName := "origin"
	flagPush := false
	flagPull := false
	flagFastforward := false
	flagForce := false
	flagRebase := "main"
	flagMerge := ""

	cmd := &cobra.Command{
		Use:   "branch:current",
		Short: "Work with the current branch",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			branchName, _ := helpers.GetCurrentBranchName()

			if flagRebase != "" {
				if flagMerge == branchName {
					fmt.Println("error: cannot rebase current branch onto itself")
					return
				}
				helpers.RunCommandOnStdout("git", "rebase", flagRebase)
			}

			if flagMerge != "" {
				if flagMerge == branchName {
					fmt.Println("error: cannot merge current branch into itself")
					return
				}
				helpers.RunCommandOnStdout("git", "merge", flagMerge, "-s", "ort")
			}

			if flagPull {
				args := []string{"pull", flagRemoteName, branchName}
				if flagFastforward {
					args = append(args, "--ff-only")
				} else {
					args = append(args, "--rebase")
				}
				fmt.Printf("args: %v\n", args)
				helpers.RunCommandOnStdout("git", args...)
			}

			if flagPush {
				args := []string{"push", flagRemoteName, branchName}
				if flagForce {
					args = append(args, "--force")
				}

				helpers.RunCommandOnStdout("git", args...)
			}

			if !flagPull && !flagPush {
				fmt.Println(branchName)
			}
		},
	}

	cmd.Flags().StringVarP(&flagRemoteName, "remote", "r", "origin", "remote name to use when pushing or pulling")
	cmd.Flags().BoolVarP(&flagPush, "push", "p", false, "push the current branch to remote")
	cmd.Flags().BoolVarP(&flagPull, "pull", "u", false, "pull the current branch from remote using rebase")
	cmd.Flags().BoolVarP(&flagFastforward, "ff", "f", false, "when pulling, use fast-forward-only")
	cmd.Flags().BoolVarP(&flagForce, "force", "F", false, "when pushing, perform a force push")
	cmd.Flags().StringVarP(&flagRebase, "rebase", "R", "", "rebase the current branch using the specified branch")
	cmd.Flags().StringVarP(&flagMerge, "merge", "M", "", "merge the specified branch into the current branch")

	rootCmd.AddCommand(cmd)
}
