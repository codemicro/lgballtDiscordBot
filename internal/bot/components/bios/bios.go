package bios

import (
	"bytes"
	_ "embed"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/markdown"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Bios struct {
	trackerLock   *sync.RWMutex
	trackedEmbeds map[string]*trackedEmbed // map of message IDs to tracked embed objects
}

type trackedEmbed struct {
	current        int
	accountId      string
	channelId      string
	bios           []db.UserBio
	timeoutAt      time.Time
	requestingUser string
	isAdmin        bool
}

const (
	maxBioFieldLen      = 1024
	nextBioReaction     = "➡️"
	previousBioReaction = "⬅️"
	bioTimeoutDuration  = time.Minute * 5
)

var (
	helpEmbed = &discordgo.MessageEmbed{
		Title: "Bios help",
	}

	systemHelpEmbed = &discordgo.MessageEmbed{
		Title: "Bios for systems",
	}
)

//go:embed biosHelp.md
var helpMarkdown string

//go:embed biosSyshelp.md
var syshelpMarkdown string

func init() {
	generateFields := func(sourceTpl string, data interface{}, emb *discordgo.MessageEmbed) error {

		tplParsed, err := template.New("").Parse(sourceTpl)
		if err != nil {
			return err
		}

		tpld := new(bytes.Buffer)
		if err = tplParsed.Execute(tpld, data); err != nil {
			return err
		}

		for _, v := range markdown.SplitByHeader(tpld.String()) {
			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: v.Title, Value: v.Content})
		}

		return nil
	}

	err := generateFields(helpMarkdown, struct {
		Fields string
	}{Fields: strings.Join(config.BioFields, "\n")}, helpEmbed)
	if err != nil {
		panic(err)
	}

	err = generateFields(syshelpMarkdown, struct {
		Ping string
	}{Ping: tools.MakePing("289130374204751873")}, systemHelpEmbed)
	if err != nil {
		panic(err)
	}
}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(Bios)
	comp.trackerLock = new(sync.RWMutex)
	comp.trackedEmbeds = make(map[string]*trackedEmbed)

	// goroutine to close tracked embeds when they timeout
	go func() {
		for {
			time.Sleep(time.Second * 30)
			comp.trackerLock.Lock()

			var toRemove []string
			for key, tracked := range comp.trackedEmbeds {
				if tracked.timeoutAt.Before(time.Now()) {
					// this embed has timed out
					toRemove = append(toRemove, key)
				}
			}

			for _, key := range toRemove {
				err := kit.Session.MessageReactionsRemoveAll(comp.trackedEmbeds[key].channelId, key)
				if err != nil {
					log.Error().Err(err).Msgf("unable to clear reactions from tracked message %s", key)
				}
				delete(comp.trackedEmbeds, key)
			}

			comp.trackerLock.Unlock()
		}
	}()

	kit.AddReaction(&route.Reaction{
		Name:  "Bio pagination reaction",
		Run:   comp.PaginationReaction,
		Event: route.ReactionAdd,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios help",
		Help:        "View information about the bios feature of the bot",
		CommandText: []string{"bio", "help"},
		Run:         comp.SingletHelp,
		Category:    meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios for systems help",
		Help:        "View information about the bios for systems feature of the bot",
		CommandText: []string{"bio", "syshelp"},
		Run:         comp.SystemHelp,
		Category:    meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios set field (systems)",
		Help:        "Set a field in a system member's bio",
		CommandText: []string{"bio", "set"},
		Arguments: []route.Argument{
			{Name: "memberId", Type: pluralkitMemberIdType{}},
			{Name: "field", Type: bioFieldType{}},
			{Name: "newValue", Type: route.RemainingString},
		},
		Run:              comp.SystemSetField,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios set field",
		Help:        "Set a field in your bio",
		CommandText: []string{"bio", "set"},
		Arguments: []route.Argument{
			{Name: "field", Type: bioFieldType{}},
			{Name: "newValue", Type: route.RemainingString},
		},
		Run:              comp.SingletSetField,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios clear field (systems)",
		Help:        "Clear a field in a system member's bio",
		CommandText: []string{"bio", "clear"},
		Arguments: []route.Argument{
			{Name: "memberId", Type: pluralkitMemberIdType{}},
			{Name: "field", Type: bioFieldType{}},
		},
		Run:              comp.SystemClearField,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios clear field",
		Help:        "Clear a field in your bio",
		CommandText: []string{"bio", "clear"},
		Arguments: []route.Argument{
			{Name: "field", Type: bioFieldType{}},
		},
		Run:              comp.SingletClearField,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios set image (systems)",
		Help:        "Set the image attached to your bio",
		CommandText: []string{"bio", "setimg"},
		Arguments: []route.Argument{
			{Name: "memberId", Type: pluralkitMemberIdType{}},
			{Name: "imageURL", Type: route.URL, Default: func(session *discordgo.Session, message *discordgo.MessageCreate) (interface{}, error) {
				return "", nil
			}},
		},
		Run:              comp.SystemSetImage,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios set image",
		Help:        "Set the image attached to your bio",
		CommandText: []string{"bio", "setimg"},
		Arguments: []route.Argument{
			{Name: "imageURL", Type: route.URL, Default: func(session *discordgo.Session, message *discordgo.MessageCreate) (interface{}, error) {
				return "", nil
			}},
		},
		Run:              comp.SingletSetImage,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios clear image (systems)",
		Help:        "Clear the image attached to a system member's bio",
		CommandText: []string{"bio", "clearimg"},
		Arguments: []route.Argument{
			{Name: "memberId", Type: pluralkitMemberIdType{}},
		},
		Run:              comp.SystemClearImage,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:             "Bios clear image",
		Help:             "Clear the image attached to your bio",
		CommandText:      []string{"bio", "clearimg"},
		Arguments:        []route.Argument{},
		Run:              comp.SingletClearImage,
		AllowOverloading: true,
		Category:         meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios import PluralKit member (systems)",
		Help:        "Import a PluralKit system member's profile to use with a bio",
		CommandText: []string{"bio", "import"},
		Arguments: []route.Argument{
			{Name: "memberId", Type: pluralkitMemberIdType{}},
		},
		Run:      comp.SystemImportMember,
		Category: meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios view",
		Help:        "View the bio(s) associated with a Discord account",
		CommandText: []string{"bio"},
		Arguments: []route.Argument{
			{Name: "user", Type: common.PingOrUserIdType, Default: func(_ *discordgo.Session, message *discordgo.MessageCreate) (interface{}, error) {
				return message.Author.ID, nil
			}},
		},
		Run:      comp.ReadBio,
		Category: meta.CategoryBios,
	})

	return nil

}
