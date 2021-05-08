package route

import "github.com/bwmarrin/discordgo"

type CommonContext struct {
	Session *discordgo.Session
	Kit     *Kit
}

func (m *CommonContext) DefaultAllowedMentions() *discordgo.MessageAllowedMentions {
	// This copy is intentional
	n := m.Kit.DefaultAllowedMentions
	return &n
}

func (m *CommonContext) SendMessageString(channelId string, content string) (*discordgo.Message, error) {

	return m.Session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content:         content,
		AllowedMentions: m.DefaultAllowedMentions(),
	})

}

func (m *CommonContext) SendMessageEmbed(channelId string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {

	return m.Session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Embed:           embed,
		AllowedMentions: m.DefaultAllowedMentions(),
	})

}
