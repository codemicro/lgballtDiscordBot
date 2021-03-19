package bios

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"strings"
	"sync"
	"time"
)

type Bios struct {
	trackerLock   *sync.RWMutex
	trackedEmbeds map[string]*trackedEmbed // map of message IDs to tracked embed objects
}

type trackedEmbed struct {
	current   int
	accountId string
	channelId string
	bios      []db.UserBio
	timeoutAt time.Time
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
		Fields: []*discordgo.MessageEmbedField{
			{Name: "What are bios?", Value: "Think of Bios like ID cards. Input responses for pre-defined fields such as a Pronouns field or a Sexuality field, and then your responses get put into a nice Bio card which can be viewed by anyone using `$bio [@username]`. \n\nIf you just want to get your own, start by filling in any of the below fields!"},
			{Name: "How do I fill out these fields?", Value: "Run `$bio set [field] [value]`. For example, `$bio set Pronouns She/Her` would set the Pronouns field of your Bio to \"She/Her\".\n\nTo remove a field, run `$bio clear [field]` with no other arguments."},
			{Name: "What fields can I fill in?", Value: fmt.Sprintf("The current fields are:\n```yml\n%s\n```", strings.Join(config.BioFields, "\n"))},
			{Name: "What about bios for systems?", Value: "I'm glad you asked! Run `$bio syshelp` for more information about bios for systems."},
			{Name: "I think a new field should be added. How can I request one?", Value: "Request new fields in <#698575463278313583>."},
			{Name: "How do I view someone's Bio without mentioning them?", Value: "User IDs can be used instead of mentioning a user. To get a User ID, first enable Developer Mode by going to User Settings, Appearance, and toggling it to on. After that, right click a username on desktop or tap the 3 dots on a profile card on mobile then click \"Copy ID\". \nNow just run `$bio [UserID]` to view their Bio. For example, `$bio 516962733497778176`"},
			{Name: "Anything else I should know?", Value: "- You don't need to wipe a field to put in new info. Just run `$bio set [field] [text]` to overwrite it.\n- If you end up in a situation where you have no fields left in your bio because you've removed them all, your entire bio is deleted."},
			{Name: "TL;DR/Commands", Value: "- View your own Bio with `$bio`, another user's with `$bio [user id or mention]`\n- Fill in a field with `$bio set [field] [text]`. Fields can be overwritten with the same command.\n- Wipe a field with `$bio clear [field]`"},
		},
	}

	systemHelpEmbed = &discordgo.MessageEmbed{
		Title: "Bios for systems",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Creating a bio", Value: "Unlike with regular bios for singlets, you need to explicitly create a bio for a system member before you can edit bio fields. This can be done using `$bio import <member ID>`, where `member ID` is the PluralKit member ID that you would like to import. If anything is in the birthday, pronouns or description fields, they will be automatically copied into the new bio."},
			{Name: "Updating bio fields", Value: "You can update and remove bio fields in a very similar fashion to bios for singlets, namely using the following two commands: `$bio clear <member ID> <field>` to remove a field and `$bio set <member ID> <field> <new contents>` to update a field's value. Field names are the same as bios for singlets, which can be found in `$bio help`."},
			{Name: "Deleting a system member bio", Value: "To delete a system member bio, simply remove every field that's present in it using `$bio clear <member ID> <field>`. This will trigger the bot to automatically delete the bio entry from the database."},
			{Name: "Viewing system member bios", Value: "System member bios can be viewed using the same command as you would use to view singlet bios, namely `$bio [ping or user ID]`. If a user account has multiple bios associated with it, it will show a carousel-type interface, allowing you to scroll between bios.\nCurrently, this is not well optimised for viewing accounts that have large numbers of associated bios, and may be a little cumbersome to use."},
			{Name: "Anything else?", Value: "Be aware that any changes you make in a bio for a system member will not be replicated in the bio that PluralKit stores for that member. Using bios for systems doesn't affect your ability to use regular bios.\nIf you have questions, encounter an issue or have a suggestion, feel free to ping Abi! (0x414b#8669)"},
		},
	}
)

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
					logging.Error(err, fmt.Sprintf("unable to clear reactions from tracked message %s", key)) // TODO: zerolog
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
		Category: meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios for systems help",
		Help:        "View information about the bios for systems feature of the bot",
		CommandText: []string{"bio", "syshelp"},
		Run:         comp.SystemHelp,
		Category: meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios set field",
		Help:        "Set a field in your bio",
		CommandText: []string{"bio", "set"},
		Arguments: []route.Argument{
			{Name: "field", Type: bioFieldType{}},
			{Name: "newValue", Type: route.RemainingString},
		},
		Run: comp.SingletSetField,
		AllowOverloading: true,
		Category: meta.CategoryBios,
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
		Run: comp.SystemSetField,
		AllowOverloading: true,
		Category: meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios clear field",
		Help:        "Clear a field in your bio",
		CommandText: []string{"bio", "clear"},
		Arguments: []route.Argument{
			{Name: "field", Type: bioFieldType{}},
		},
		Run: comp.SingletClearField,
		AllowOverloading: true,
		Category: meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios clear field (systems)",
		Help:        "Clear a field in a system member's bio",
		CommandText: []string{"bio", "clear"},
		Arguments: []route.Argument{
			{Name: "memberId", Type: pluralkitMemberIdType{}},
			{Name: "field", Type: bioFieldType{}},
		},
		Run: comp.SystemClearField,
		AllowOverloading: true,
		Category: meta.CategoryBios,
	})

	kit.AddCommand(&route.Command{
		Name:        "Bios import PluralKit member (systems)",
		Help:        "Import a PluralKit system member's profile to use with a bio",
		CommandText: []string{"bio", "import"},
		Arguments: []route.Argument{
			{Name: "memberId", Type: pluralkitMemberIdType{}},
		},
		Run: comp.SystemImportMember,
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
		Run: comp.ReadBio,
		Category: meta.CategoryBios,
	})

	return nil

}
