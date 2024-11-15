package cmd

import (
	"fmt"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

func init() {
	var flagAutoPull bool = false

	var checkoutCmd = &cobra.Command{
		Use:     "checkout",
		Aliases: []string{"co"},
		Short:   "Checks out the specified branch",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("error: branch name required")
				return
			}

			if err := helpers.RunCommandOnStdout("git", "checkout", args[0]); err != nil {
				return
			}

			if flagAutoPull {
				currentBranch, _ := helpers.GetCurrentBranchName()

				if currentBranch != args[0] {
					fmt.Println("error: failed to switch branches")
					return
				}

				if err := helpers.RunCommandOnStdout("git", "pull", "origin", args[0], "--rebase"); err != nil {
					return
				}
			}
		},
	}

	rootCmd.AddCommand(checkoutCmd)
	checkoutCmd.Flags().BoolVarP(&flagAutoPull, "pull", "p", false, "Automatically pull origin after checkout")
}
