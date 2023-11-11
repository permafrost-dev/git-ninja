package cmd

import (
	"fmt"
	"os/exec"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch:current:push",
		Short: "Push the current branch to origin",
		Run: func(c *cobra.Command, args []string) {
			branch, err := helpers.GetCurrentBranchName()

			if err != nil {
				fmt.Println(err)
				return
			}

			exec.Command("git", "push", "origin", branch).Run()
		},
	})
}
