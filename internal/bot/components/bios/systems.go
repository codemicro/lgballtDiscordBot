package bios

import (
	"context"
	"errors"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
	"github.com/skwair/harmony"
	"regexp"
	"strings"
)

func (b *Bios) HelpSystem(_ []string, m *harmony.Message) error {
	_, err := b.b.SendMessage(m.ChannelID, "TODO") // TODO
	return err
}

func (b *Bios) SetFieldSystem(command []string, m *harmony.Message) error {
	// Syntax: <member ID> <field name> <value>

	newValue := strings.Join(command[2:], " ")
	bdt := new(db.UserBio)
	bdt.UserId = m.Author.ID
	bdt.SysMemberID = command[0]

	return b.setBioField(bdt, command[1], newValue, true, m)
}

func (b *Bios) ClearFieldSystem(command []string, m *harmony.Message) error {
	// Syntax: <member ID> <field name>

	bdt := new(db.UserBio)
	bdt.UserId = m.Author.ID
	bdt.SysMemberID = command[0]

	return b.clearBioField(bdt, command[1], m)
}

var sysmateDetectionBio = regexp.MustCompile(`(?m)[a-zA-Z]{5}`)

func (b *Bios) ImportSystemMember(command []string, m *harmony.Message) error {
	// Syntax: <member ID>

	if len(command) < 1 {
		return nil
	}

	memberId := strings.ToLower(command[0])

	if !sysmateDetectionBio.MatchString(memberId) {
		// poorly formatted sysmate bio
		_, err := b.b.SendMessage(m.ChannelID, "This member ID is invalid - it must be in the format `abcde`.")
		return err
	}

	// check to see if member already imported
	accBios, err := db.GetBiosForAccount(m.Author.ID)
	if err != nil {
		return err
	}

	for _, x := range accBios {
		if strings.EqualFold(x.SysMemberID, memberId) {
			_, err := b.b.SendMessage(m.ChannelID, "This member ID has already been imported!")
			return err
		}
	}

	// get bio
	pkMember, err := pluralkit.MemberByMemberId(memberId)
	if err != nil {
		if errors.Is(err, pluralkit.ErrorMemberNotFound) {
			_, err := b.b.SendMessage(m.ChannelID, "This member ID is not recognised by the PluralKit API (is" +
				" this a valid member ID?)\nThe PluralKit API returned a HTTP 404.")
			return err
		}
		return err
	}

	var otherText, pronounsText string

	if pkMember.Description != "" {
		otherText += pkMember.Description
	}

	if pkMember.Birthday != "" {
		if pkMember.Description != "" {
			otherText += "\n\n"
		}
		otherText += "Birthday: " + pkMember.Birthday
	}

	if pkMember.Pronouns != "" {
		pronounsText = pkMember.Pronouns
	}

	if otherText == "" && pronounsText == "" {
		otherText = "Placeholder value"
	}

	bdt := new(db.UserBio)
	bdt.UserId = m.Author.ID
	bdt.SysMemberID = command[0]
	bdt.BioData = make(map[string]string)

	if otherText != "" {
		bdt.BioData["Other"] = otherText
	}

	if pronounsText != "" {
		bdt.BioData["Pronouns"] = pronounsText
	}

	err = bdt.Create()
	if err != nil {
		return err
	}

	// react to message with a check mark to signify it worked
	err = b.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
	return err
}
