package bios

import (
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/skwair/harmony/embed"
	"strings"
	"time"
)

type Bios struct {
	b *core.Bot
	data biosData
}

func New(bot *core.Bot) (*Bios, error) {
	b := new(Bios)
	b.b = bot

	dt, err := loadBiosFile()
	if err != nil {
		return nil, err
	}
	b.data = dt

	go func() {
		for {
			time.Sleep(time.Second * 5)
			err := saveBiosFile(b.data)
			if err != nil {
				logging.Error(err)
			}
		}

	}()

	return b, nil
}

var biosHelpEmbed = embed.New().
	Title("Bio help/FAQ").
	Fields(
		embed.NewField().Name("What are bios?").Value("Think of Bios like ID cards. Input responses for pre-defined fields such as a Pronouns field or a Sexuality field, and then your responses get put into a nice Bio card which can be viewed by anyone using `$bio [@username]`. \n\nIf you just want to get your own, run `$bio` with no other arguments.").Build(),
		embed.NewField().Name("How do I input responses for these fields?").Value("Run `$bio [field] [value]`. For example, `$bio Pronouns She/Her` would set the Pronouns field of your Bio to \"She/Her\".\n\nTo remove a field, run `$bio [field]` with no other arguments.").Build(),
		embed.NewField().Name("What fields can I fill in?").Value("The current fields are:\n```yml\nNames\nNicknames\nSysmates\nPronouns\nGender\nSexuality\nRomanticism\nOrientation\nIdentities\nOther```").Build(),
		embed.NewField().Name("I think a new field should be added. How can I request one?").Value("Request new fields in <#698575463278313583>.").Build(),
		embed.NewField().Name("How do I view someone's Bio without mentioning them?").Value("User IDs can be used instead of mentioning a user. To get a User ID, first enable Developer Mode by going to User Settings, Appearance, and toggling it to on. After that, right click a username on desktop or tap the 3 dots on a profile card on mobile then click \"Copy ID\". \nNow just run `$bio [UserID]` to view their Bio. For example, `$bio 516962733497778176`").Build(),
		embed.NewField().Name("Anything else I should know?").Value("- You don't need to wipe a field to put in new info. Just run `$bio [field] [text]` to overwrite it.").Build(),
		embed.NewField().Name("TL;DR/Commands").Value("- View your own Bio with `$bio`, another user's with `$bio [user id or mention`\n- Fill in a field with `$bio [field] [text]`. Fields can be overwritten with the same command.\n- Wipe a field with `$bio [field]`\n- View only a specific field on a Bio with `$bio [user id or mention] [field]`").Build(),
	).Build()

// ValidateFieldName performs a case insensitive compare of the provided field name and those used in the data file
// If a match is found, the properly capitalised version is returned.
func (b *Bios) ValidateFieldName(inputName string) (properFieldName string, found bool) {
	for _, f := range b.data.Fields {
		if strings.EqualFold(f, inputName) {
			found = true
			properFieldName = f
			break
		}
	}
	return
}