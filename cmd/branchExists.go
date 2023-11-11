package cmd

import (
	"fmt"
	"os"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch:exists [name]",
		Short: "Check if the given branch name exists",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("error: branch name required")
				return
			}

			if exists, _ := helpers.BranchExists(args[0]); !exists {
				os.Exit(1)
			}

			os.Exit(0)
		},
	})
}
