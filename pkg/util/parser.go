package util

import (
	"regexp"
	"strings"
)

// Command represents a parsed chat-ops command
type Command struct {
	Name string
	Args []string
}

var commandRegex = regexp.MustCompile(`^/(\w+)(?:\s+(.*))?$`)

// ParseCommand extracts a command from a comment body
// Returns nil if no valid command found
func ParseCommand(body string) *Command {
	// Split into lines and find first line starting with /
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "/") {
			continue
		}

		matches := commandRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		cmd := &Command{
			Name: matches[1],
		}

		if len(matches) > 2 && matches[2] != "" {
			// Split arguments by whitespace
			cmd.Args = strings.Fields(matches[2])
		}

		return cmd
	}

	return nil
}
