package bios

import (
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
)

var sysmateNotRegisteredMessage = "This member is not registered to your Discord account with the bot. Please import " +
	"this member using `$bio import %s`"

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
			err = ctx.SendErrorMessage(fmt.Sprintf(sysmateNotRegisteredMessage, bdt.SysMemberID))
			return err
		}
		err = bdt.Create()
	} else {
		err = bdt.Save()
	}

	if err != nil {
		return err
	}

	// react to message with a check mark to signify it worked
	err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
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

	if len(bdt.BioData) == 0 && bdt.ImageURL != "" {
		// There's no fields and no image left in the bio, so we shall delete it
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
		err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bios) setBioImage(bdt *db.UserBio, imageURL string, isSysmate bool, ctx *route.MessageContext) error {

	hasBio, err := bdt.Populate()
	if err != nil {
		return err
	}

	var selectedImage string
	var sendWarning bool
	if imageURL == "" {
		if len(ctx.Message.Attachments) == 0 {
			return ctx.SendErrorMessage("no image provided - please attach a file to this message or provide an image URL")
		}
		selectedImage = ctx.Message.Attachments[0].URL
		sendWarning = true
	} else {
		selectedImage = imageURL
	}

	bdt.ImageURL = selectedImage

	if !hasBio {
		if isSysmate {
			err = ctx.SendErrorMessage(fmt.Sprintf(sysmateNotRegisteredMessage, bdt.SysMemberID))
			return err
		}
		err = bdt.Create()
	} else {
		err = bdt.Save()
	}

	if err != nil {
		return err
	}

	// react to message with a check mark to signify it worked
	err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
	if err != nil {
		return err
	}

	if sendWarning {
		_, err := ctx.SendMessageString(ctx.Message.ChannelID, "Image successfully added\nPlease note that if you" +
			" remove the message containing the image, the image will no longer show on your bio.")
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bios) clearBioImage(bdt *db.UserBio, ctx *route.MessageContext) error {

	hasBio, err := bdt.Populate()
	if err != nil {
		return err
	}

	if !hasBio {
		err = ctx.SendErrorMessage("You have not created a bio, hence there's no image to clear.")
		return err
	}

	bdt.ImageURL = ""

	if len(bdt.BioData) == 0 && bdt.ImageURL != "" {
		// There's no fields and no image left in the bio, so we shall delete it
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
		err = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, v)
		if err != nil {
			return err
		}
	}

	return nil
}