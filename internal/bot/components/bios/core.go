package bios

import (
	"context"
	"fmt"
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
		fmt.Println("created")
		err = bdt.Create()
	} else {
		fmt.Println("saved")
		err = bdt.Save()
	}

	if err != nil {
		return err
	}

	// react to message with a check mark to signify it worked
	err = b.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
	return err
}