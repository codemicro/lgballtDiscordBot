package bios

import (
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
	"strings"
)

func (*Bios) SystemHelp(ctx *route.MessageContext) error {
	_, err := ctx.SendMessageEmbed(ctx.Message.ChannelID, systemHelpEmbed)
	return err
}

func (b *Bios) SystemSetField(ctx *route.MessageContext) error {

	memberId := ctx.Arguments["memberId"].(string)
	newValue := ctx.Arguments["newValue"].(string)
	field := ctx.Arguments["field"].(string)

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID
	bdt.SysMemberID = memberId

	return b.setBioField(bdt, field, newValue, true, ctx)
}

func (b *Bios) SystemClearField(ctx *route.MessageContext) error {

	memberId := ctx.Arguments["memberId"].(string)
	field := ctx.Arguments["field"].(string)

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID
	bdt.SysMemberID = memberId

	return b.clearBioField(bdt, field, ctx)
}

func (b *Bios) SystemSetImage(ctx *route.MessageContext) error {

	memberId := ctx.Arguments["memberId"].(string)
	imageURL := ctx.Arguments["imageURL"].(string)

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID
	bdt.SysMemberID = memberId

	return b.setBioImage(bdt, imageURL, true, ctx)
}

func (b *Bios) SystemClearImage(ctx *route.MessageContext) error {

	memberId := ctx.Arguments["memberId"].(string)

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID
	bdt.SysMemberID = memberId

	return b.clearBioImage(bdt, ctx)
}

func (b *Bios) SystemImportMember(ctx *route.MessageContext) error {

	memberId := ctx.Arguments["memberId"].(string)

	// check to see if member already imported
	accBios, err := db.GetBiosForAccount(ctx.Message.Author.ID)
	if err != nil {
		return err
	}

	for _, x := range accBios {
		if strings.EqualFold(x.SysMemberID, memberId) {
			err = ctx.SendErrorMessage("This member ID has already been imported!")
			return err
		}
	}

	// check to see if account has a system
	systemInfo, err := pluralkit.SystemById(ctx.Message.Author.ID)
	if err != nil {
		if e, ok := err.(*pluralkit.Error); ok {
			if e.Code == pluralkit.ErrorCodeSystemNotFound {
				return ctx.SendErrorMessage("Your Discord account has no PluralKit systems associated with it.")
			}
		}
		return err
	}

	// check system has the specified member ID as a listed member

	systemMembers, err := pluralkit.MembersBySystemId(systemInfo.Id)
	if err != nil {

		if e, ok := err.(*pluralkit.Error); ok {
			if e.Code == pluralkit.ErrorCodeUnauthorizedMemberList {
				return ctx.SendErrorMessage("Your system has the member list set to **private**. Please set " +
					"this to public (`pk;s privacy members public`) and try again (you can set the member list to " +
					"private again afterwards)")
			}
		}

		return err
	}

	var pkMember *pluralkit.Member
	for _, sysm := range systemMembers {
		if strings.EqualFold(sysm.Id, memberId) {
			pkMember = sysm
			break
		}
	}

	if pkMember == nil {
		return ctx.SendErrorMessage(fmt.Sprintf("Your system has has no member with the given ID. If you're " +
			"sure there's a registered member with this ID, make sure the member visibility privacy level is set " +
			"to **public** (`pk;m %s privacy visibility public`). ", memberId))
	}

	// make bio

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
	bdt.UserId = ctx.Message.Author.ID
	bdt.SysMemberID = memberId
	bdt.SystemID = systemInfo.UUID
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
	err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
	return err
}
