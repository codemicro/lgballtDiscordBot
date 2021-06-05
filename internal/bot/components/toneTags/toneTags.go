package toneTags

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

type ToneTags struct{}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(ToneTags)

	kit.AddCommand(&route.Command{
		Name:        "Tone tag lookup",
		Help:        "Look up the definition of a tone tag",
		CommandText: []string{"toneTag", "lookup"},
		Arguments: []route.Argument{
			{Name: "tag", Type: route.String},
		},
		Run:      comp.Lookup,
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "List all tone tags",
		Help:        "List all tone tags known to the bot",
		CommandText: []string{"toneTag", "list"},
		Run:         comp.List,
		Category:    meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Register new tone tag",
		Help:        "Register a new tone tag description",
		CommandText: []string{"toneTag", "new"},
		Arguments: []route.Argument{
			{Name: "tag", Type: route.String},
			{Name: "description", Type: route.RemainingString},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRole),
		},
		Run:      comp.Create,
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Delete a tone tag",
		CommandText: []string{"toneTag", "delete"},
		Arguments: []route.Argument{
			{Name: "tag", Type: route.String},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRole),
		},
		Run:      comp.Delete,
		Category: meta.CategoryMisc,
	})

	return nil
}
