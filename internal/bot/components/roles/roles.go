package roles

import "github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"

const partyRole = "698570587567685703"

type Roles struct {
	b    *core.Bot
}

func New(bot *core.Bot) (*Roles, error) {
	b := new(Roles)
	b.b = bot
	return b, nil
}