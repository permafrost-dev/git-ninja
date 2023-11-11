package utils

import (
	"os/exec"
	"strings"
)

func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(result)), nil
}
