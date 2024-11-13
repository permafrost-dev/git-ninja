package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

// getAvailableBranches fetches and returns a map of all available branches in the repository
func getAvailableBranches() map[string]bool {
	cmd := exec.Command("git", "branch", "--list")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return make(map[string]bool)
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

	return branchMap
}

// branchExists checks if a branch exists in the cached branch list
func branchExists(branch string, availableBranches map[string]bool) bool {
	_, exists := availableBranches[branch]
	return exists
}

func getGitReflogLines() ([]string, error) {
	output, err := helpers.RunCommandBuffered("git", "reflog", "show", "--pretty=format:'%at ~ %gs ~ %gd'", "--date=relative")
	if err != nil {
		return make([]string, 0), err
	}

	return strings.Split(output, "\n"), nil
}

func branchMatchesFilterIgnoreFlag(branch string) bool {
	if len(flagFilterIgnore) > 0 {
		if matched, _ := regexp.MatchString(flagFilterIgnore, branch); matched {
			return true
		}
	}

	return false
}

func parseTimestampIntoTime(timestamp string) time.Time {
	seconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(seconds, 0)
}

var flagCount int = 10
var flagFilterIgnore string = ""

type BranchInfo struct {
	Name      string
	Age       string
	Timestamp time.Time
}

var listRecentBranchesCmd = &cobra.Command{
	Use:   "branch:recent [--count|-c <count>]",
	Short: "Show recently checked out branch names",
	Run: func(c *cobra.Command, args []string) {
		existingBranches := getAvailableBranches()
		lines, _ := getGitReflogLines()

		lineRegex := regexp.MustCompile(`([0-9]+) ~ (checkout):.+ ([^~]+) ~ HEAD@{(.*)}`)
		seen := make(map[string]bool)
		count := 0

		for _, line := range lines {
			matches := lineRegex.FindStringSubmatch(line)
			if len(matches) < 4 {
				continue
			}

			info := BranchInfo{
				Name:      strings.TrimSpace(matches[3]),
				Age:       strings.TrimSpace(matches[4]),
				Timestamp: parseTimestampIntoTime(matches[1]),
			}

			// exclude branches that are not in the list of current branches, i.e. branches that have been deleted
			if !branchExists(info.Name, existingBranches) {
				continue
			}

			if branchMatchesFilterIgnoreFlag(info.Name) {
				continue
			}

			if !seen[info.Name] {
				seen[info.Name] = true

				fmt.Printf("  \033[33m%-15s \033[37;1m %s\033[0m\n", info.Age, info.Name)

				if count += 1; count >= flagCount {
					break
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listRecentBranchesCmd)

	listRecentBranchesCmd.Flags().IntVarP(&flagCount, "count", "c", 10, "Limit the number of branches to display")
	listRecentBranchesCmd.Flags().StringVarP(&flagFilterIgnore, "exclude", "e", "", "Exclude branches that match the provided regex")
}
