package messageTools

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

type MessageTools struct{}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(MessageTools)

	kit.AddCommand(&route.Command{
		Name:        "Send message",
		Help:        "Send a message in a specific channel",
		CommandText: []string{"m", "send"},
		Arguments: []route.Argument{
			{Name: "channelId", Type: route.DiscordSnowflake},
			{Name: "message", Type: route.RemainingString},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRoles...),
		},
		Run: comp.Send,
		Category: meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Edit message",
		Help:        "Edit a message the bot has sent",
		CommandText: []string{"m", "edit"},
		Arguments: []route.Argument{
			{Name: "messageLink", Type: route.URL},
			{Name: "newContent", Type: route.RemainingString},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRoles...),
		},
		Run: comp.Edit,
		Category: meta.CategoryAdminTools,
	})

	return nil

}
