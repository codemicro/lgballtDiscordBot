package muteme

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

var muteMeText = "This will mute you until %s and cannot be undone. Are you sure?"

type MuteMe struct{}

func Init(kit *route.Kit, st *state.State) error {

	comp := new(MuteMe)
	go comp.startMuteRemovalWorker(kit.Session, st)

	kit.AddCommand(&route.Command{
		Name:        "Mute me",
		Help:        "Procrastinating? Want to lock yourself out of the server to stop getting distracted? Mute yourself with this handy command! :D",
		CommandText: []string{"muteme"},
		Arguments: []route.Argument{
			{Name: "duration", Type: route.Duration},
		},
		Run: comp.Trigger,
	})

	return nil

}
