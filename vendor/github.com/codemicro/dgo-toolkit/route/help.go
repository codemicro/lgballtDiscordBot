package route

import "strings"

// CommandInfo contains textual information about a command
type CommandInfo struct {
	Name            string
	Description     string
	CommandText     string
	HasRestrictions bool
	Category        uint
	Arguments       []*ArgumentInfo
}

// ArgumentInfo contains textual information about a command argument
type ArgumentInfo struct {
	Name       string
	Type       string
	Usage      string
	HasDefault bool
}

// GetCommandInfo returns information about all commands in the kit (intended for use in help commands)
func (b *Kit) GetCommandInfo() []*CommandInfo {
	var n []*CommandInfo
	for _, cmd := range b.commandSet {

		if cmd.Invisible {
			continue
		}

		var args []*ArgumentInfo
		for _, ag := range cmd.Arguments {
			args = append(args, &ArgumentInfo{
				Name:       ag.Name,
				Type:       ag.Type.Name(),
				Usage:      ag.Type.Help(ag.Name),
				HasDefault: ag.Default != nil,
			})
		}

		n = append(n, &CommandInfo{
			Name:            cmd.Name,
			Description:     cmd.Help,
			CommandText:     strings.Join(cmd.CommandText, " "),
			Arguments:       args,
			HasRestrictions: len(cmd.Restrictions) != 0,
			Category:        cmd.Category,
		})
	}

	return n
}

// GetCommandInfo returns information about all commands in the kit (intended for use in help commands)
func (m *MessageContext) GetCommandInfo() []*CommandInfo {
	return m.Kit.GetCommandInfo()
}
