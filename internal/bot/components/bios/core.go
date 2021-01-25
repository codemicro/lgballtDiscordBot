package bios

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"strings"
	"sync"
)

func (b *Bios) setBioField(bdt *db.UserBio, rawFieldName, newValue string, m *harmony.Message) error {
	fieldName, validFieldName := b.validateFieldName(rawFieldName)
	if !validFieldName {
		_, err := b.b.SendMessage(m.ChannelID, "That's not a valid field name! Choose from one of the "+
			"following: "+strings.Join(b.data.Fields, ", "))
		return err
	}

	if len(newValue) > maxBioFieldLen {
		_, err := b.b.SendMessage(m.ChannelID, "Sorry - the new text you have entered is too long (this is a "+
			"Discord limitation). Please limit each field of your bio to `1024` characters.")
		return err
	}

	hasBio, err := bdt.Populate()
	if err != nil {
		return err
	}

	if !hasBio {
		bdt.BioData = make(map[string]string)
	}

	bdt.BioData[fieldName] = newValue

	if !hasBio {
		err = bdt.Create()
	} else {
		err = bdt.Save()
	}

	if err != nil {
		return err
	}

	// react to message with a check mark to signify it worked
	err = b.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
	return err
}

func (b *Bios) clearBioField(bdt *db.UserBio, rawFieldName string, m *harmony.Message) error {
	fieldName, validFieldName := b.validateFieldName(rawFieldName)
	if !validFieldName {
		_, err := b.b.SendMessage(m.ChannelID, "That's not a valid field name! Choose from one of the "+
			"following: "+strings.Join(b.data.Fields, ", "))

		return err
	}

	hasBio, err := bdt.Populate()
	if err != nil {
		return err
	}

	if !hasBio {  // This theoretically will never happen because of the MID check on the route phase, but I'm leaving a
		// check here anyway
		_, err := b.b.SendMessage(m.ChannelID, "You have not created a bio, hence there is nothing to delete anything from.")
		return err
	}

	delete(bdt.BioData, fieldName)

	if len(bdt.BioData) == 0 {
		// There are no fields left in the bio, so we shall delete it
		err = bdt.Delete()
	} else {
		// Else save as normal
		err = bdt.Save()
	}

	if err != nil {
		return err
	}

	// react to message with a check mark to signify it worked
	for _, v := range []string{"🗑", "✅"} {
		err := b.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, v)
		if err != nil {
			return err
		}
	}

	return nil
}

type pluralityInfo struct {
	CurrentNumber int
	TotalCount int
}

type nameDriver interface {
	Name() (string, error)
	Avatar() (string, error)

	SysMemberId() string
	HasMultiple() bool
	CurrentAndTotalCount() (int, int)
}

type systemName struct {
	memberId string
	plurality *pluralityInfo

	member   *pluralkit.Member
	once     *sync.Once
}

func newSystemName(memberId string, info *pluralityInfo) *systemName {
	return &systemName{
		memberId: memberId,
		plurality: info,
	}
}

func (sn *systemName) fetchInformation() error {
	if sn.once == nil {
		sn.once = new(sync.Once)
	}

	var err error
	sn.once.Do(func() {
		var member *pluralkit.Member
		member, err = pluralkit.MemberByMemberId(sn.memberId)
		if err != nil {
			return
		}
		sn.member = member
	})
	return err
}

func (sn *systemName) Name() (string, error) {
	if err := sn.fetchInformation(); err != nil {
		return "", err
	}

	var name string
	if sn.member.Nickname == "" {
		name = sn.member.Name
	} else {
		name = sn.member.Nickname
	}

	return name, nil
}

func (sn *systemName) Avatar() (string, error) {
	if err := sn.fetchInformation(); err != nil {
		return "", err
	}
	return sn.member.Avatar, nil
}

func (sn *systemName) SysMemberId() string {
	return sn.memberId
}

func (sn *systemName) HasMultiple() bool {
	return sn.plurality != nil
}

func (sn *systemName) CurrentAndTotalCount() (int, int) {
	if sn.plurality != nil {
		return sn.plurality.CurrentNumber, sn.plurality.TotalCount
	}
	return 0, 0
}

type accountName struct {
	accountId string
	guildId string

	name string
	avatar string

	plurality *pluralityInfo

	user *harmony.User
	bot *core.Bot

	once *sync.Once
}

func newAccountName(accountId, guildId string, info *pluralityInfo, bot *core.Bot) *accountName {
	return &accountName{
		accountId: accountId,
		guildId:   guildId,
		plurality: info,
		bot: bot,
	}
}

func (an *accountName) fetchInformation() error {
	if an.once == nil {
		an.once = new(sync.Once)
	}

	var err error
	an.once.Do(func() {
		var name string
		var user *harmony.User
		name, user, err = an.bot.GetNickname(an.accountId, an.guildId)
		if err != nil {
			return
		}
		an.name = name
		an.avatar = user.AvatarURL()
	})
	return err
}

func (an *accountName) Name() (string, error) {
	if err := an.fetchInformation(); err != nil {
		return "", err
	}
	return an.name, nil
}

func (an *accountName) Avatar() (string, error) {
	if err := an.fetchInformation(); err != nil {
		return "", err
	}
	return an.avatar, nil
}

func (an *accountName) SysMemberId() string {
	return ""
}

func (an *accountName) HasMultiple() bool {
	return an.plurality != nil
}

func (an *accountName) CurrentAndTotalCount() (int, int) {
	if an.plurality != nil {
		return an.plurality.CurrentNumber, an.plurality.TotalCount
	}
	return 0, 0
}

// formBioEmbed creates an embed object based on a user's bio data
func (b *Bios) formBioEmbed(nd nameDriver, bioData map[string]string) (*embed.Embed, error) {

	name, err := nd.Name()
	if err != nil {
		return nil, err
	}

	avatar, err := nd.Avatar()
	if err != nil {
		return nil, err
	}

	var footerText string

	if nd.HasMultiple() {
		footerText += "This account has multiple bios associated with it.\n"
		curr, total := nd.CurrentAndTotalCount()
		footerText += fmt.Sprintf("Currently viewing No. %d of %d", curr, total)

		if nd.SysMemberId() != "" {
			footerText += fmt.Sprintf("\nPluralKit member ID: %s", nd.SysMemberId())
		}
	}

	e := embed.New()
	e.Thumbnail(embed.NewThumbnail(avatar))
	e.Title(fmt.Sprintf("%s's bio", name))
	e.Footer(embed.NewFooter().Text(footerText).Build())

	var fields []*embed.Field
	for _, category := range b.data.Fields {
		fVal, ok := bioData[category]
		if ok {
			fields = append(fields, embed.NewField().Name(category).Value(fVal).Build())
		}
	}

	e.Fields(fields...)

	return e.Build(), nil
}