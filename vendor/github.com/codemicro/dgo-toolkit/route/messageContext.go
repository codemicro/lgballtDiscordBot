package route

import "github.com/bwmarrin/discordgo"

type MessageContext struct {
	*CommonContext
	Message   *discordgo.MessageCreate
	Arguments map[string]interface{}
	Raw       string
	// Command must not be modified
	Command   *Command
}

// SendErrorMessage sends a Discord message containing error information using kit.UserErrorFunc
func (m *MessageContext) SendErrorMessage(problem string) error {
	_, err := m.SendMessageString(m.Message.ChannelID, m.Kit.UserErrorFunc(problem))
	return err
}
