package tools

import "strings"

func GetCommand(commandString, prefix string) []string {
	command := strings.Split(strings.Trim(commandString, " "), " ")

	if len(command) >= 1 {
		if command[0] == prefix {
			command = command[1:]
		}
	}

	if len(command) == 1 && command[0] == "" {
		command = []string{}
	}
	return command
}
