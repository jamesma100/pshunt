package parser

import (
	"bufio"
	"strings"
)

func ParseList(contents string) []string {
	scanner := bufio.NewScanner(strings.NewReader(contents))

	var psList []string

	for scanner.Scan() {
		line := scanner.Text()
		psList = append(psList, line)
	}
	return psList
}
