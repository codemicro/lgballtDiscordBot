package bios

import (
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
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

	// TODO: This is temporary, and needs properly updating. Perhaps also moving to core?

	bios, err := db.GetBiosForAccount(id)
	if err != nil {
		return err
	}

	if len(bios) == 0 {
		_, err := b.b.SendMessage(m.ChannelID, "This user hasn't created a bio, or just plain doesn't exist.")
		if err != nil {
			return err
		}
	}

	// Show first bio

	if len(bios) == 1 {
		// Found a bio, now to form an embed
		e, err := b.formBioEmbed(newAccountName(id, m.GuildID, nil, b.b), bios[0].BioData)
		if err != nil {
			return err
		}

		_, err = b.b.SendEmbed(m.ChannelID, e)
		if err != nil {
			return err
		}
	} else {

		totalBios := len(bios)

		for i, bio := range bios {
			// Found a bio, now to form an embed

			plurality := &pluralityInfo{
				CurrentNumber: i+1,
				TotalCount:    totalBios,
			}

			var nd nameDriver
			if bio.SysMemberID != "" {
				nd = newSystemName(bio.SysMemberID, plurality)
			} else {
				nd = newAccountName(id, m.GuildID, plurality, b.b)
			}

			e, err := b.formBioEmbed(nd, bio.BioData)
			if err != nil {
				return err
			}

			_, err = b.b.SendEmbed(m.ChannelID, e)
			if err != nil {
				return err
			}
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

