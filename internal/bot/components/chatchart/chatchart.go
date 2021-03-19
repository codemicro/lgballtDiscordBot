package chatchart

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

type ChatChart struct {
	queue chan collectionIntent
}

type collectionIntent struct {
	ChannelId string
	Ctx       *route.MessageContext
}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(ChatChart)
	comp.queue = make(chan collectionIntent, 1024*1024)

	go func() {
		for x := range comp.queue {
			comp.collectMessages(x)
		}
	}()

	kit.AddCommand(&route.Command{
		Name:        "Trigger chatchart",
		Help:        "Add a channel to the chatchart collection queue",
		CommandText: []string{"chatchart"},
		Arguments: []route.Argument{
			{Name: "channel", Type: route.ChannelMention},
		},
		Run: comp.Trigger,
		Category: meta.CategoryFun,
	})

	return nil
}
