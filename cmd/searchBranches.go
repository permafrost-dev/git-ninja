package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type branchInfo struct {
	branch string
	count  int
	latest time.Time
}

func getAllBranchesWithAges(availableBranches map[string]bool) ([]branchInfo, error) {
	cmd := exec.Command("git", "reflog", "show", "--pretty=format:%gs~%ci")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing git command:", err)
		return nil, err
	}

	// Convert command output to string and split by new lines
	output := out.String()
	lines := strings.Split(output, "\n")

	// Prepare regular expressions to match 'checkout:' and extract the branch name and ISO date
	checkoutRegexp := regexp.MustCompile(`^checkout: moving from [^ ]+ to ([^ ]+)~(.*)$`)

	branchData := make(map[string]branchInfo)
	branchDataArray := make([]branchInfo, 0)

	// Iterate through each line and apply the regex filters
	for _, line := range lines {
		if matches := checkoutRegexp.FindStringSubmatch(line); len(matches) == 3 {
			branch := matches[1]
			isoDate := strings.TrimSpace(matches[2])

			checkoutDate, err := time.Parse("2006-01-02 15:04:05 -0700", isoDate)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				continue
			}

			if branchExists(branch, availableBranches) {
				info := branchData[branch]
				info.count++

				// Update latest checkout date
				if info.latest.IsZero() || checkoutDate.After(info.latest) {
					info.latest = checkoutDate
				}

				branchData[branch] = info
			}
		}
	}

	for branch, info := range branchData {
		branchDataArray = append(branchDataArray, branchInfo{branch: branch, count: info.count, latest: info.latest})
	}

	for branch := range availableBranches {
		if _, ok := branchData[branch]; !ok {
			branchDataArray = append(branchDataArray, branchInfo{branch: branch, count: 0, latest: time.Time{}})
		}
	}

	return branchDataArray, nil
}

var flagRegex bool = false

var searchBranchesCmd = &cobra.Command{
	Use:   "branch:search [--regex|-r] <substring-or-regex>",
	Short: "Search branch names for matching substrings",
	Long:  `Searches branch names for matching substrings and displays a list of matching branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No search string provided.")
			return
		}
		searchFor := args[0]

		availableBranches, err := getAvailableBranches()

		if err != nil {
			fmt.Printf("Error retrieving branches: %v\n", err)
			return
		}

		availableBranchAgesSorted, err := getAllBranchesWithAges(availableBranches)

		if err != nil {
			fmt.Printf("Error retrieving branch ages: %v\n", err)
			return
		}

		// Search for branches that contain the search string
		var matches []branchInfo
		for _, branchData := range availableBranchAgesSorted {
			if !flagRegex && strings.Contains(branchData.branch, searchFor) {
				matches = append(matches, branchData)
			}

			if flagRegex {
				if matched, _ := regexp.MatchString(searchFor, branchData.branch); matched {
					matches = append(matches, branchData)
				}
			}
		}

		// Print out the matching branch names
		if len(matches) == 0 {
			fmt.Println("No matching branches found.")
			return
		}

		sort.Slice(matches, func(i, j int) bool {
			return matches[i].latest.After(matches[j].latest)
		})

		for _, branch := range matches {
			var branchAge string = "never"

			if !branch.latest.IsZero() {
				branchAge = getRelativeTime(branch.latest)
			}

			fmt.Printf("  \033[33m%-16s \033[37;1m %s\033[0m\n", branchAge, branch.branch)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchBranchesCmd)

	searchBranchesCmd.Flags().BoolVarP(&flagRegex, "regex", "r", false, "Search using a regular expression pattern")
}
