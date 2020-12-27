package chatchart

import (
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/skwair/harmony"
)

type ChatChart struct {
	b     *core.Bot
	queue chan collectionIntent
}

func New(bot *core.Bot) (*ChatChart, error) {
	b := new(ChatChart)
	b.b = bot
	b.queue = make(chan collectionIntent, 1024*1024)

	go func() {
		for {
			b.collectMessages(<-b.queue)
		}
	}()

	return b, nil
}

type collectionIntent struct {
	ChannelId string
	Message   *harmony.Message
}
