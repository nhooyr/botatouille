package parser

import "strings"

const commandPrefix = "!"

func IsCommand(message string) (cmdLine string, ok bool) {
	if strings.HasPrefix(message, commandPrefix) {
		return strings.TrimPrefix(message, commandPrefix), true
	}
	return "", false
}

func NextCommand(cmdLine string) (cmd string, rest string) {
	s := strings.SplitN(cmdLine, " ", 2)
	if len(s) == 1 {
		return cmdLine, ""
	}
	return s[0], s[1]
}