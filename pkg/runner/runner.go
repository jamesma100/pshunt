package runner

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func GetPsList() string {
	cmd := exec.Command("ps", "ax")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(stdout)
}

func KillPs(pid string) string {
	cmd := exec.Command("kill", "-9", pid)
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		return string(stderr)
	}
	return fmt.Sprintf("Killed process with pid %s", pid)
}

func Grep(needle string, haystack []string) []string {
	var filtered []string

	for _, item := range haystack {
		if strings.Contains(item, needle) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
