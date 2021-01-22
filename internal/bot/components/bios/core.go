package bios

import (
	"context"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/skwair/harmony"
	"strings"
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