package githelpers

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetLastCheckedoutBranchName() (string, error) {
	var out bytes.Buffer

	cmd := exec.Command("git", "reflog", "show", "--pretty=format:%gs", "--date=relative")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Scan through each line of output
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains 'checkout:'
		if strings.Contains(line, "checkout:") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				return fields[3], nil // Return the fourth word (branch name)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no checkout entries found")
}
