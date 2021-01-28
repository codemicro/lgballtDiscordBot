package bios

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/skwair/harmony/embed"
	"strings"
	"sync"
	"time"
)

const (
	maxBioFieldLen      = 1024
	nextBioReaction     = "➡️"
	previousBioReaction = "⬅️"
	bioTimeoutDuration  = time.Minute * 5
)

type Bios struct {
	b             *core.Bot
	data          biosData
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

var (
	biosHelpEmbed *embed.Embed
)

const (
	systemBiosPt1 = "**__Bios for systems__**\nBios for systems works similarly to bios divided by singlets, except they have close PluralKit integration.\n\n**Creating a bio for a system member**\nFor this, you will need the member ID of the system member you want to add a bio for (a name will **not** work).\n```\n$bio import <member ID>\n```\nThis will first verify that the specified member ID exists and then import the pronouns, birthday and description fields of a PluralKit member information card into a new bio.\nIf any of your member privacy settings restrict access to these fields, they will not be imported and instead left blank.\nCurrently, you can only import one system member at a time. If you'd like to import a large number, get in touch with Abi.\n\n**Updating a system member's bio**\nYou can update and delete the bio of a system member with the PluralKit member ID. The commands to do this behave identically to those for singlets with the addition of an option.\nTo delete a field:\n```\n$bio <member ID> <field>\n```\nTo set the value of a field:\n```\n$bio <member ID> <field> <content>\n```\nUse the same fields as you would for a singlet bio.\n\n**Deleting a system member's bio**\nIf all fields are removed from a system member's bio, the record will be deleted from the bot's database. If you want to add a bio for this system member again, simply follow the same process as before."
	systemBiosPt2 = "**Viewing system member bios**\nSystem member bios can be viewed using the same command as you would use to view singlet bios, namely\n```\n$bio [ping or user ID]\n```\nIf a user account has multiple bios associated with it, it will show a carousel-type interface, allowing you to scroll between user bios.\nCurrently, this is not especially optimised for viewing bios of accounts that have large numbers of member bios. I've got an idea of how to remedy this, but in order to get bios for systems released as soon as possible, it's not a feature yet. Hang tight, I'm hoping to get it implemented within a month of the v3 bot release.\n\n**Anything else?**\nBe aware that any changes you make in a bio for a system member will not be replicated in the bio that PluralKit stores for that member. This is for two reasons, mainly that PluralKit bios are limited to a much shorter content length than bios here, and updating PluralKit changes would require you to provide your PluralKit API key to the bot.\nUsing bios for systems doesn't affect your ability to use regular bios. If you'd like to go back to using only bios for singlets, simply delete every system member's bio.\nThis is also a new thing - if you have questions, encounter an issue or have a suggestion, feel free to ping Abi! (0x414b#8669)"
)

func New(bot *core.Bot) (*Bios, error) {
	b := new(Bios)
	b.b = bot
	b.trackerLock = new(sync.RWMutex)
	b.trackedEmbeds = make(map[string]*trackedEmbed)

	// goroutine to close tracked embeds when they timeout
	go func() {
		for {
			time.Sleep(time.Second * 30)
			b.trackerLock.Lock()

			var toRemove []string
			for key, tracked := range b.trackedEmbeds {
				if tracked.timeoutAt.Before(time.Now()) {
					// this embed has timed out
					toRemove = append(toRemove, key)
				}
			}

			for _, key := range toRemove {
				err := b.b.Client.Channel(b.trackedEmbeds[key].channelId).RemoveAllReactions(context.Background(), key)
				if err != nil {
					logging.Error(err, fmt.Sprintf("unable to clear reactions from tracked message %s", key))
				}
				delete(b.trackedEmbeds, key)
			}

			b.trackerLock.Unlock()
		}
	}()

	dt, err := loadBiosFile()
	if err != nil {
		return nil, err
	}
	b.data = dt

	biosHelpEmbed = embed.New().
		Title("Bio help/FAQ").
		Fields(
			embed.NewField().Name("What are bios?").Value("Think of Bios like ID cards. Input responses for pre-defined fields such as a Pronouns field or a Sexuality field, and then your responses get put into a nice Bio card which can be viewed by anyone using `$bio [@username]`. \n\nIf you just want to get your own, start by filling in any of the below fields!").Build(),
			embed.NewField().Name("How do I input responses for these fields?").Value("Run `$bio [field] [value]`. For example, `$bio Pronouns She/Her` would set the Pronouns field of your Bio to \"She/Her\".\n\nTo remove a field, run `$bio [field]` with no other arguments.").Build(),
			embed.NewField().Name("What fields can I fill in?").Value(fmt.Sprintf("The current fields are:\n```yml\n%s\n```", strings.Join(b.data.Fields, "\n"))).Build(),
			embed.NewField().Name("What about bios for systems?").Value("I'm glad you asked! Run `$bio syshelp` for more information about bios for systems.").Build(),
			embed.NewField().Name("I think a new field should be added. How can I request one?").Value("Request new fields in <#698575463278313583>.").Build(),
			embed.NewField().Name("How do I view someone's Bio without mentioning them?").Value("User IDs can be used instead of mentioning a user. To get a User ID, first enable Developer Mode by going to User Settings, Appearance, and toggling it to on. After that, right click a username on desktop or tap the 3 dots on a profile card on mobile then click \"Copy ID\". \nNow just run `$bio [UserID]` to view their Bio. For example, `$bio 516962733497778176`").Build(),
			embed.NewField().Name("Anything else I should know?").Value("- You don't need to wipe a field to put in new info. Just run `$bio [field] [text]` to overwrite it.\n- If you end up in a situation where you have no fields left in your bio because you've removed them all, your entire bio is deleted.").Build(),
			embed.NewField().Name("TL;DR/Commands").Value("- View your own Bio with `$bio`, another user's with `$bio [user id or mention]`\n- Fill in a field with `$bio [field] [text]`. Fields can be overwritten with the same command.\n- Wipe a field with `$bio [field]`\n- View only a specific field on a Bio with `$bio [user id or mention] [field]`").Build(),
		).Build()

	return b, nil
}

// validateFieldName performs a case insensitive compare of the provided field name and those used in the data file
// If a match is found, the properly capitalised version is returned.
func (b *Bios) validateFieldName(inputName string) (properFieldName string, found bool) {
	for _, f := range b.data.Fields {
		if strings.EqualFold(f, inputName) {
			found = true
			properFieldName = f
			break
		}
	}
	return
}
