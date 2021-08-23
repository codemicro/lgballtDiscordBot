package roles

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

type Roles struct{}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(Roles)

	kit.AddCommand(&route.Command{
		Name:        "Add reaction role",
		Help:        "Track an emoji on a specific message to use it as a reaction role",
		CommandText: []string{"roles", "track"},
		Arguments: []route.Argument{
			{Name: "messageLink", Type: route.URL},
			{Name: "emoji", Type: route.String},
			{Name: "roleName", Type: route.RemainingString},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRoles...),
		},
		Run: comp.Track,
		Category: meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Remove reaction role",
		Help:        "Untrack an emoji on a specific message",
		CommandText: []string{"roles", "untrack"},
		Arguments: []route.Argument{
			{Name: "messageLink", Type: route.URL},
			{Name: "emoji", Type: route.String},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRoles...),
		},
		Run: comp.Untrack,
		Category: meta.CategoryAdminTools,
	})

	kit.AddReaction(&route.Reaction{
		Name:  "Reaction roles: add",
		Run:   comp.ReactionAdd,
		Event: route.ReactionAdd,
	}, &route.Reaction{
		Name:  "Reaction roles: remove",
		Run:   comp.ReactionRemove,
		Event: route.ReactionRemove,
	})

	return nil

}
