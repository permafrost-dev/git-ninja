package helpers

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/permafrost-dev/git-ninja/app/utils"
)

func RunCommandOnStdout(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

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

func GetCurrentBranchName() (string, error) {
	result, err := utils.RunCommand("git", "branch", "--show-current")
	result = strings.TrimSpace(result)

	if strings.Contains(result, " ") {
		return "", fmt.Errorf(result)
	}

	return result, err
}

func BranchExists(name string) (bool, error) {
	var out bytes.Buffer

	cmd := exec.Command("git", "branch", "--list")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return false, err
	}

	// Scan through each line of output
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "* ")

		if strings.EqualFold(line, name) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
