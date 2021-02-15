package messagetools

import "github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"

type MessageTools struct {
	b *core.Bot
}

func New(bot *core.Bot) (*MessageTools, error) {
	b := new(MessageTools)
	b.b = bot

	return b, nil
}
