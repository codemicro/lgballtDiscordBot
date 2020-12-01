package info

import (
	"github.com/codemicro/lgballtDiscordBot/internal/bot"
)

type info struct {
	b *bot.Bot
	commandPrefix string
}

func Register(bot *bot.Bot, commandName string) error {
	b := new(info)
	b.b = bot
	b.commandPrefix = bot.Prefix + commandName

	bot.Client.OnMessageCreate(b.onMessageCreate)

	return nil
}
