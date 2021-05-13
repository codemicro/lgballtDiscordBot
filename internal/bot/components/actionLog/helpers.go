package actionLog

import "github.com/bwmarrin/discordgo"

const (
	unknownAuthor = "`<unknown>`"
	cannotGetContent = "unable to retrieve content"
)

func getAuthorMention(m *discordgo.Message) (string, bool) {
	if m != nil && m.Author != nil {
		return m.Author.Mention(), true
	}
	return unknownAuthor, false
}

func getAuthorUsername(m *discordgo.Message) (string, bool) {
	if m != nil && m.Author != nil {
		return m.Author.String(), true
	}
	return unknownAuthor, false
}

func getContent(m *discordgo.Message) (string, bool) {
	if m != nil {
		if x := m.ContentWithMentionsReplaced(); len(x) != 0 {
			return x, true
		}
	}
	return cannotGetContent, false
}

func wasBot(m *discordgo.Message) bool {
	return m != nil && m.Author != nil && m.Author.Bot
}