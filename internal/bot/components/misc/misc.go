package misc

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
)

type Misc struct{}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(Misc)

	kit.AddCommand(&route.Command{
		Name:        "Avatar",
		Help:        "Get an enlarged version of a user's avatar",
		CommandText: []string{"avatar"},
		Arguments: []route.Argument{
			{Name: "userId", Type: route.String, Default: func(_ *discordgo.Session, message *discordgo.MessageCreate) (interface{}, error) {
				return message.Author.ID, nil
			}},
		},
		Run: comp.Avatar,
	})

	kit.AddCommand(&route.Command{
		Name:        "Emoji",
		Help:        "Get an enlarged version of a custom emoji",
		CommandText: []string{"emoji"},
		Arguments: []route.Argument{
			{Name: "emoji", Type: route.String},
		},
		Run: comp.Emoji,
	})

	kit.AddCommand(&route.Command{
		Name:        "Steal emojis",
		Help:        "Commit theivery and steal someone else's custom emojis",
		CommandText: []string{"steal"},
		Arguments: []route.Argument{
			{Name: "messageLink", Type: route.URL},
		},
		Run: comp.StealEmojis,
	})

	kit.AddCommand(&route.Command{
		Name:        "Help",
		Help:        "Show command information",
		CommandText: []string{"help"},
		Run:         comp.Help,
		Invisible:   true,
	})

	kit.AddCommand(&route.Command{
		Name:        "Request listener",
		Help:        "Ping someone with the listener role to listen to you in a venting channel",
		CommandText: []string{"listener"},
		Restrictions: []route.CommandRestriction{
			func(_ *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
				return tools.IsStringInSlice(message.ChannelID, config.Listeners.AllowedChannels), nil
			},
		},
		Run: comp.ListenToMe,
	})

	return nil

}
