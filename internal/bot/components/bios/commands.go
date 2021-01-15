package bios

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"strings"
)

// Help sends the bios help embed message
func (b *Bios) Help(_ []string, m *harmony.Message) error {
	_, err := b.b.SendEmbed(m.ChannelID, biosHelpEmbed)
	return err
}

// ReadBio runs the bio read command
func (b *Bios) ReadBio(command []string, m *harmony.Message) error {
	// Syntax: <user ID>

	var id string
	if len(command) >= 1 {
		// If there's a ping as the argument, use the ID from that. Else, just use the plain argument
		id = command[0]
		if x, match := tools.ParsePing(id); match {
			id = x
		}
	} else {
		// Since no ID argument is provided, assume it's that of the message author
		id = m.Author.ID
	}

	bdt := new(db.UserBio)
	ok, err := bdt.Populate(id)
	if err != nil {
		return err
	}
	bioData := bdt.BioData

	if !ok {
		// No bio for that user
		_, err := b.b.SendMessage(m.ChannelID, "This user hasn't created a bio, or just plain doesn't exist.")
		if err != nil {
			return err
		}
	} else {
		// Found a bio, now to form an embed
		e, err := b.formBioEmbed(id, m.GuildID, bioData)
		if err != nil {
			return err
		}

		_, err = b.b.SendEmbed(m.ChannelID, e)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetField performs the bio field set command
func (b *Bios) SetField(command []string, m *harmony.Message) error {
	// Syntax: <field name> <value>

	properFieldName, validFieldName := b.ValidateFieldName(command[0])

	if !validFieldName {
		_, err := b.b.SendMessage(m.ChannelID, "That's not a valid field name! Choose from one of the "+
			"following: "+strings.Join(b.data.Fields, ", "))

		return err
	}

	newValue := strings.Join(command[1:], " ")

	if len(newValue) > maxBioFieldLen {
		_, err := b.b.SendMessage(m.ChannelID, "Sorry - the new text you have entered is too long (this is a "+
			"Discord limitation). Please limit each field of your bio to `1024` characters.")
		return err
	}

	bdt := new(db.UserBio)
	hasBio, err := bdt.Populate(m.Author.ID)
	if err != nil {
		return err
	}
	bioData := bdt.BioData

	if !hasBio {
		bioData = make(map[string]string)
		bioData[properFieldName] = newValue
	} else {
		bioData[properFieldName] = newValue
	}

	bdt.BioData = bioData

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

// ClearField runs the bio field clear command
func (b *Bios) ClearField(command []string, m *harmony.Message) error {
	// Syntax: <field name>

	properFieldName, validFieldName := b.ValidateFieldName(command[0])

	if !validFieldName {
		_, err := b.b.SendMessage(m.ChannelID, "That's not a valid bio field.")
		return err
	}

	bdt := new(db.UserBio)
	hasBio, err := bdt.Populate(m.Author.ID)
	if err != nil {
		return err
	}

	if !hasBio {
		_, err := b.b.SendMessage(m.ChannelID, "You have not created a bio, hence there is nothing to delete anything from.")
		return err
	}

	delete(bdt.BioData, properFieldName)

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

// formBioEmbed creates an embed object based on a user's bio data
func (b *Bios) formBioEmbed(uid, guildId string, bioData map[string]string) (*embed.Embed, error) {

	var name string
	var avatar string

	name, user, err := b.b.GetNickname(uid, guildId)
	if err != nil {
		return nil, err
	}
	avatar = user.AvatarURL()

	e := embed.New()
	e.Thumbnail(embed.NewThumbnail(avatar))
	e.Title(fmt.Sprintf("%s's bio", name))

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
