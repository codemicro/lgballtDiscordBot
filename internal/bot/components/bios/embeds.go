package bios

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
	"strconv"
	"sync"
)

type pluralityInfo struct {
	CurrentNumber int
	TotalCount    int
}

type nameDriver interface {
	Name() (string, error)
	Avatar() (string, error)
	Colour() (int, error)

	SysMemberId() string
	HasMultiple() bool
	CurrentAndTotalCount() (int, int)
}

type systemName struct {
	memberId  string
	colour    int
	plurality *pluralityInfo

	member *pluralkit.Member
	once   *sync.Once
}

func newSystemName(memberId string, info *pluralityInfo) *systemName {
	return &systemName{
		memberId:  memberId,
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

		if member.Colour != "" {
			var col int64
			col, err = strconv.ParseInt(member.Colour, 16, 32)
			(fmt.Sprint(col))
			if err != nil {
				return
			}
			sn.colour = int(col)
		}

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

func (sn *systemName) Colour() (int, error) {
	if err := sn.fetchInformation(); err != nil {
		return 0, err
	}
	return sn.colour, nil
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
	guildId   string

	name   string
	avatar string
	colour int

	plurality *pluralityInfo

	user    *discordgo.User
	session *discordgo.Session

	once *sync.Once
}

func newAccountName(accountId, guildId string, info *pluralityInfo, session *discordgo.Session) *accountName {
	return &accountName{
		accountId: accountId,
		guildId:   guildId,
		plurality: info,
		session:   session,
	}
}

func (an *accountName) fetchInformation() error {
	if an.once == nil {
		an.once = new(sync.Once)
	}

	var err error
	an.once.Do(func() {
		var name string
		var user *discordgo.User
		name, user, err = common.GetNickname(an.session, an.accountId, an.guildId)
		if err != nil {
			return
		}
		an.name = name
		an.avatar = user.AvatarURL("")

		member, memberErr := an.session.GuildMember(an.guildId, an.accountId)
		if memberErr == nil {
			if len(member.Roles) != 0 {

				var roles []*discordgo.Role
				roles, err = an.session.GuildRoles(an.guildId)
				if err != nil {
					return
				}

				var role *discordgo.Role
				for _, memberRole := range member.Roles {
					for _, r := range roles {
						if r.ID == memberRole && r.Color != 0 { // first role with a colour set
							role = r
							break
						}
					}
				}


				if role == nil {
					return
				}

				an.colour = role.Color
			}
		}
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

func (an *accountName) Colour() (int, error) {
	if err := an.fetchInformation(); err != nil {
		return 0, err
	}
	return an.colour, nil
}

func (an *accountName) CurrentAndTotalCount() (int, int) {
	if an.plurality != nil {
		return an.plurality.CurrentNumber, an.plurality.TotalCount
	}
	return 0, 0
}

// formBioEmbed creates an embed object based on a user's bio data
func (b *Bios) formBioEmbed(nd nameDriver, bioData map[string]string) (*discordgo.MessageEmbed, error) {

	name, err := nd.Name()
	var avatar string
	var colour int
	// if any error is returned from the Name function, it will also affect the Avatar and Colour function, but will
	// not be returned there.
	if err != nil {
		logging.Warn(err.Error())
		name = "<unknown name>"
	} else {
		avatar, _ = nd.Avatar()
		colour, _ = nd.Colour()
	}

	var footerText string

	if nd.HasMultiple() {

		if errors.Is(err, pluralkit.ErrorMemberNotFound) {
			footerText += "⚠ This member appears to have been deleted from PluralKit ⚠\n"
		}

		footerText += "This account has multiple bios associated with it.\n"
		curr, total := nd.CurrentAndTotalCount()
		footerText += fmt.Sprintf("Currently viewing No. %d of %d", curr, total)

		if nd.SysMemberId() != "" {
			footerText += fmt.Sprintf("\nPluralKit member ID: %s", nd.SysMemberId())
		}
	}

	var fields []*discordgo.MessageEmbedField
	for _, category := range config.BioFields {
		fVal, ok := bioData[category]
		if ok {
			fields = append(fields, &discordgo.MessageEmbedField{Name: category, Value: fVal})
		}
	}

	e := discordgo.MessageEmbed{
		Title:     fmt.Sprintf("%s's bio", name),
		Footer:    &discordgo.MessageEmbedFooter{Text: footerText},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: avatar},
		Fields:    fields,
		Color:     colour,
	}

	return &e, nil
}
