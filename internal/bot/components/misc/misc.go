package misc

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

type Misc struct {}

func Init(kit *route.Kit, state *state.State) error {

	comp := new(Misc)

	kit.AddCommand(&route.Command{
		Name:         "Avatar",
		Help:         "Get an enlarged version of a user's avatar",
		CommandText:  []string{"avatar"},
		Arguments:    []route.Argument{
			{Name: "userId", Type: route.String, Default: func(_ *discordgo.Session, message *discordgo.MessageCreate) (interface{}, error) {
				return message.Author.ID, nil
			}},
		},
		Restrictions: nil,
		Run:          comp.Avatar,
	})

	return nil

}
