package misc

import "github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"

type Misc struct {
	b *core.Bot
}

func New(bot *core.Bot) (*Misc, error) {
	b := new(Misc)
	b.b = bot

	go b.startMuteRemovalWorker()

	return b, nil
}
