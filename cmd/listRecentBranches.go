package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

// getAvailableBranches fetches and returns a map of all available branches in the repository
func getAvailableBranches() (map[string]bool, error) {
	cmd := exec.Command("git", "branch", "--list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Convert the command output to a map for quick lookup
	branches := strings.Split(out.String(), "\n")
	branchMap := make(map[string]bool)
	for _, branch := range branches {
		branch = strings.TrimSpace(strings.TrimPrefix(branch, "*"))
		if branch != "" {
			branchMap[branch] = true
		}
	}

	return branchMap, nil
}

// branchExists checks if a branch exists in the cached branch list
func branchExists(branch string, availableBranches map[string]bool) bool {
	_, exists := availableBranches[branch]
	return exists
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch:recent",
		Short: "Show recently checked out branch names",
		Run: func(c *cobra.Command, args []string) {

			availableBranches, err := helpers.GetAvailableBranchesMap()
			if err != nil {
				fmt.Println("Error fetching branches:", err)
				return
			}

			cmd := exec.Command("git", "reflog", "show", "--pretty=format:%gs ~ %gd", "--date=relative")
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error executing git command:", err)
				return
			}

			// Convert command output to string and split by new lines
			output := out.String()
			lines := strings.Split(output, "\n")

			// Prepare regular expressions to match 'checkout:' and to extract the relevant parts
			checkoutRegexp := regexp.MustCompile(`checkout:`)
			extractRegexp := regexp.MustCompile(`([^ ]+) ~ (.*)`)

			// Keep track of seen branches to avoid duplicates
			seen := make(map[string]bool)
			count := 0

			// Iterate through each line and apply the regex filters
			for _, line := range lines {
				if !checkoutRegexp.MatchString(line) {
					continue
				}

				matches := extractRegexp.FindStringSubmatch(line)
				if len(matches) < 3 {
					continue
				}

				branch := matches[1]
				head := strings.TrimSpace(matches[2])

				if !seen[branch] && branchExists(branch, availableBranches) {
					seen[branch] = true
					// Format the output

					age := strings.TrimSuffix(head, "}")
					age = strings.Replace(age, "HEAD@{", "", -1)
					fmt.Printf("  \033[33m%s: \033[37;1m %s\033[0m\n", age, branch)
					count++
				}

				// Limit output to 20 lines
				if count >= 20 {
					break
				}
			}
		},
	})
}
