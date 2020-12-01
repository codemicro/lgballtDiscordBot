package info

import (
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
)

type Info struct {
	b *core.Bot
}

func New(bot *core.Bot) (*Info, error) {
	b := new(Info)
	b.b = bot

	return b, nil
}
