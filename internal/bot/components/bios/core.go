package bios

import (
	"errors"
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
)

func (*Bios) setBioField(bdt *db.UserBio, fieldName, newValue string, isSysmate bool, ctx *route.MessageContext) error {

	if len(newValue) > maxBioFieldLen {
		err := ctx.SendErrorMessage("Sorry - the new text you have entered is too long (this is a Discord" +
			" limitation). Please limit each field of your bio to `1024` characters.")
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
		if isSysmate {
			err = ctx.SendErrorMessage(fmt.Sprintf("This member is not registered to your Discord account with"+
				" the bot. Please import this member using `$bio import %s`", bdt.SysMemberID))
			return err
		}
		err = bdt.Create()
	} else {
		err = bdt.Save()
	}

	if err != nil {
		return err
	}

	// When a message is proxied by PluralKit, this will apply the reactions to the correct message (and not the one
	// that's just been/just about to be deleted)
	// TODO: Turn this into a standalone function
	targetMessageID := ctx.Message.ID
	pkMsg, err := pluralkit.MessageById(targetMessageID)
	if err != nil && !errors.Is(err, pluralkit.ErrorMessageNotFound) {
		logging.Warn(err.Error())
	} else if pkMsg != nil {
		targetMessageID = pkMsg.Id
	}

	// react to message with a check mark to signify it worked
	err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, targetMessageID, "✅")
	return err
}

func (b *Bios) clearBioField(bdt *db.UserBio, fieldName string, ctx *route.MessageContext) error {

	hasBio, err := bdt.Populate()
	if err != nil {
		return err
	}

	if !hasBio {
		err = ctx.SendErrorMessage("You have not created a bio, hence there is nothing to delete anything " +
			"from.")
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

	// When a message is proxied by PluralKit, this will apply the reactions to the correct message (and not the one
	// that's just been/just about to be deleted)
	// TODO: Turn this into a standalone function
	targetMessageID := ctx.Message.ID
	pkMsg, err := pluralkit.MessageById(targetMessageID)
	if err != nil && !errors.Is(err, pluralkit.ErrorMessageNotFound) {
		logging.Warn(err.Error())
	} else if pkMsg != nil {
		targetMessageID = pkMsg.Id
	}

	// react to message with a check mark to signify it worked
	for _, v := range []string{"🗑", "✅"} {
		err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, targetMessageID, v)
		if err != nil {
			return err
		}
	}

	return nil
}
