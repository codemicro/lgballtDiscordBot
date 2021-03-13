package info

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

type Info struct {}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(Info)

	kit.AddCommand(&route.Command{
		Name:         "Ping",
		Help:         "Ping the bot and get the current heartbeat latency",
		CommandText:  []string{"info", "ping"},
		Run:          comp.Ping,
	})

	kit.AddCommand(&route.Command{
		Name:         "Info",
		Help:         "Get information about the bot",
		CommandText:  []string{"info"},
		Run:          comp.Info,
	})

	return nil
}