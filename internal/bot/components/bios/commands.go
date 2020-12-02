package bios

import (
	"context"
	"fmt"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"regexp"
	"strings"
)

var idFromPingRegex = regexp.MustCompile(`(?m)<@!(.+)>`)

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
		if v := idFromPingRegex.FindStringSubmatch(command[0]); len(v) > 1 {
			id = v[1]
		}
	} else {
		// Since no ID argument is provided, assume it's that of the message author
		id = m.Author.ID
	}

	b.data.Lock.RLock()
	bioData, ok := b.data.UserBios[id]
	b.data.Lock.RUnlock()

	if !ok {
		// No bio for that user
		_, err := b.b.SendMessage(m.ChannelID, "This user hasn't created a bio, or just plain doesn't exist.")
		if err != nil {
			return err
		}
	} else {
		// Found a bio, now to form an embed
		e, err := b.formBioEmbed(id, bioData)
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

	b.data.Lock.RLock()
	bioData, hasBio := b.data.UserBios[m.Author.ID]
	b.data.Lock.RUnlock()

	newValue := strings.Join(command[1:], " ")

	if !hasBio {
		bioData = make(map[string]string)
		bioData[properFieldName] = newValue
	} else {
		bioData[properFieldName] = newValue
	}

	b.data.Lock.Lock()
	b.data.UserBios[m.Author.ID] = bioData
	b.data.Lock.Unlock()

	// react to message with a check mark to signify it worked
	err := b.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
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

	b.data.Lock.RLock()
	bioData, hasBio := b.data.UserBios[m.Author.ID]
	b.data.Lock.RUnlock()
	if !hasBio {
		_, err := b.b.SendMessage(m.ChannelID, "You have not created a bio, hence there is nothing to delete anything from.")
		return err
	}

	delete(bioData, properFieldName)

	b.data.Lock.Lock()
	b.data.UserBios[m.Author.ID] = bioData
	b.data.Lock.Unlock()

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
func (b *Bios) formBioEmbed(uid string, bioData map[string]string) (*embed.Embed, error) {

	// Get user that owns this bio
	bioUser, err := b.b.Client.User(context.Background(), uid)
	if err != nil {
		return &embed.Embed{}, err
	}

	e := embed.New()
	e.Thumbnail(embed.NewThumbnail(bioUser.AvatarURL()))
	e.Title(fmt.Sprintf("%s's bio", bioUser.Username))

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
