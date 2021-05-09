package actionLog

import "github.com/bwmarrin/discordgo"

func getAuthorMention(m *discordgo.Message) string {
	if m != nil && m.Author != nil {
		return m.Author.Mention()
	}
	return "`<unknown>`"
}

func getAuthorUsername(m *discordgo.Message) string {
	if m != nil && m.Author != nil {
		return m.Author.String()
	}
	return "`<unknown>`"
}

func getContent(m *discordgo.Message) string {
	if m != nil {
		if x := m.ContentWithMentionsReplaced(); len(x) != 0 {
			return x
		}
	}
	return "unable to retrieve content"
}

func wasBot(m *discordgo.Message) bool {
	return m != nil && m.Author != nil && m.Author.Bot
}