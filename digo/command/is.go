package command

import "strings"

const prefix = "!"

// Is checks if the message is a command.
// If it is, it returns the command line with the prefix stripped.
func Is(message string) (cmdLine string, ok bool) {
	if strings.HasPrefix(message, prefix) {
		return strings.TrimPrefix(message, prefix), true
	}
	return "", false
}