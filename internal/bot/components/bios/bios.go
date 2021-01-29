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
	systemBiosText = "__**Bios for systems**__\n\n**Creating a bio**\nUnlike with regular bios for singlets, you need to explicitly create a bio for a system member before you can edit bio fields. This can be done using `$bio import <member ID>`, where `member ID` is the PluralKit member ID that you would like to import. If anything is in the birthday, pronouns or description fields, they will be automatically copied into the new bio.\n\n**Updating bio fields**\nYou can update and remove bio fields in a very similar fashion to bios for singlets, namely using the following two commands: `$bio <member ID> <field>` to remove a field and `$bio <member ID> <field> <new contents>` to update a field's value. Field names are the same as bios for singlets, which can be found in `$bio help`.\n\n**Deleting a system member bio**\nTo delete a system member bio, simply remove every field that's present in it using `$bio <member ID> <field>`. This will trigger the bot to automatically delete the bio entry from the database.\n\n**Viewing system member bios**\nSystem member bios can be viewed using the same command as you would use to view singlet bios, namely `$bio [ping or user ID]`. If a user account has multiple bios associated with it, it will show a carousel-type interface, allowing you to scroll between bios.\nCurrently, this is not well optimised for viewing accounts that have large numbers of associated bios. A solution to this should (hopefully) be coming within a month or so of the v3 release of the bot.\n\n**Anything else?**\nBe aware that any changes you make in a bio for a system member will not be replicated in the bio that PluralKit stores for that member. Using bios for systems doesn't affect your ability to use regular bios.\nThis is also a new thing - if you have questions, encounter an issue or have a suggestion, feel free to ping Abi! (0x414b#8669)"
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
