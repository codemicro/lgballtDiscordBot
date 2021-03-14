package pressf

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"sync"
)

type PressF struct {
	active map[string]*activePressF
	mux    *sync.RWMutex
}

const regionalFEmoji = "🇫"

func Init(kit *route.Kit, _ *state.State) error {

	comp := &PressF{
		active: make(map[string]*activePressF),
		mux:    new(sync.RWMutex),
	}

	kit.AddCommand(&route.Command{
		Name:        "Press F",
		Help:        "Press F to pay your respects, and let everyone else do it too",
		CommandText: []string{"pressf"},
		Arguments: []route.Argument{
			{Name: "thing", Type: route.RemainingString},
		},
		Run: comp.Trigger,
	})

	kit.AddReaction(&route.Reaction{
		Name:  "Press F reaction",
		Run:   comp.Reaction,
		Event: route.ReactionAdd,
	})

	return nil

}
