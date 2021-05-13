package misc

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"os"
	"sync"
	"time"
)

type Misc struct {

	helpEmbeds []*discordgo.MessageEmbed
	helpEmbedOnce *sync.Once

}

func Init(kit *route.Kit, runState *state.State) error {

	comp := new(Misc)
	comp.helpEmbedOnce = new(sync.Once)

	kit.AddCommand(&route.Command{
		Name:        "Avatar",
		Help:        "Get an enlarged version of a user's avatar",
		CommandText: []string{"avatar"},
		Arguments: []route.Argument{
			{Name: "user", Type: common.PingOrUserIdType, Default: func(_ *discordgo.Session, message *discordgo.MessageCreate) (interface{}, error) {
				return message.Author.ID, nil
			}},
		},
		Run: comp.Avatar,
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Emoji",
		Help:        "Get an enlarged version of a custom emoji",
		CommandText: []string{"emoji"},
		Arguments: []route.Argument{
			{Name: "emoji", Type: route.String},
		},
		Run: comp.Emoji,
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Steal emojis",
		Help:        "Commit theivery and steal someone else's custom emojis",
		CommandText: []string{"steal"},
		Arguments: []route.Argument{
			{Name: "messageLink", Type: route.URL},
		},
		Run: comp.StealEmojis,
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Help",
		Help:        "Show command information",
		CommandText: []string{"help"},
		Run:         comp.Help,
		Invisible:   true,
		Category: meta.CategoryMisc,
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
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Forget me",
		Help:        "Delete all data associated with your Discord account from the bot database",
		CommandText: []string{"forgetme"},
		Run: comp.ForgetMe,
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Dump your data",
		Help:        "Get a dump of all the data associated with your Discord user ID in the bot database",
		CommandText: []string{"mydata"},
		Run: comp.GetMyData,
		Category: meta.CategoryMisc,
	})

	kit.AddCommand(&route.Command{
		Name:        "Reset days since last incident",
		Help:        "Reset the number of days since the last incident to zero",
		CommandText: []string{"incidents", "reset"},
		Restrictions: []route.CommandRestriction{route.RestrictionByRole(config.AdminRole)},
		Run: comp.ResetSinceLastIncident,
		Category: meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Days since last incident",
		Help:        "Show the number of days since the last incident",
		CommandText: []string{"incidents"},
		Restrictions: []route.CommandRestriction{route.RestrictionByRole(config.AdminRole)},
		Run: comp.SinceLastIncident,
		Category: meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Shut down bot",
		CommandText: []string{"shutdown"},
		Restrictions: []route.CommandRestriction{route.RestrictionByRole(config.AdminRole)},
		Run: func(ctx *route.MessageContext) error {
			logging.Info(fmt.Sprintf("Shutting down by request of %s %s", ctx.Message.Author.ID, ctx.Message.Author.String()))
			_, _ = ctx.SendMessageString(ctx.Message.ChannelID, "***oh no what no A-[the earth stops rotating]***")
			runState.TriggerShutdown()
			timedOut := runState.WaitUntilAllComplete(time.Second * 10)
			if timedOut {
				_, _ = ctx.SendMessageString(ctx.Message.ChannelID, "Shutdown timeout exceeded, forcibly shutting down")
				logging.Warn("Shutdown timeout exceeded, forcibly shutting down")
			} else {
				_, _ = ctx.SendMessageString(ctx.Message.ChannelID, "Shutdown successful, goodbye.")
			}
			fmt.Println("Shutting down...")
			os.Exit(0)
			return nil
		},
		Category: meta.CategoryAdminTools,
	})
	
	return nil

}
