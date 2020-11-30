package tools

import "strings"

func GetCommand(commandString, prefix string) []string {
	fullCommand := strings.TrimPrefix(commandString, prefix)
	fullCommand = strings.Trim(fullCommand, " ")
	command := strings.Split(fullCommand, " ")
	if len(command) == 1 && command[0] == "" {
		command = []string{}
	}
	return command
}
