package bios

import (
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

	// TODO: carousel for systems

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
	bdt.UserId = id
	ok, err := bdt.Populate()
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

	newValue := strings.Join(command[1:], " ")
	bdt := new(db.UserBio)
	bdt.UserId = m.Author.ID

	return b.setBioField(bdt, command[0], newValue, m)
}

// ClearField runs the bio field clear command
func (b *Bios) ClearField(command []string, m *harmony.Message) error {
	// Syntax: <field name>

	bdt := new(db.UserBio)
	bdt.UserId = m.Author.ID

	return b.clearBioField(bdt, command[0], m)
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
