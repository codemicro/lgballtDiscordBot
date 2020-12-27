package bios

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/skwair/harmony"
	"strings"
)

// AdminReadBio runs the bio read command
func (b *Bios) AdminReadRawBio(command []string, m *harmony.Message) error {
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

	bdt := new(db.UserBio)
	ok, err := bdt.Populate(id)
	if err != nil {
		return err
	}

	if !ok {
		// No bio for that user
		_, err := b.b.SendMessage(m.ChannelID, "This user hasn't created a bio, or just plain doesn't exist.")
		if err != nil {
			return err
		}
	} else {
		_, err := b.b.SendMessage(m.ChannelID, fmt.Sprintf("User ID: `%s`\n```%s```", id, bdt.RawBioData))
		if err != nil {
			return err
		}
	}

	return nil
}

// AdminSetRawBio performs the bio field set command
func (b *Bios) AdminSetRawBio(command []string, m *harmony.Message) error {
	// Syntax: <user ID> <value>

	// If there's a ping as the argument, use the ID from that. Else, just use the plain argument
	id := command[0]
	if v := idFromPingRegex.FindStringSubmatch(command[0]); len(v) > 1 {
		id = v[1]
	}

	bdt := new(db.UserBio)
	hasBio, err := bdt.Populate(id)
	if err != nil {
		return err
	}

	completedEmojis := []string{"✅"}

	cmdContent := strings.Join(command[1:], " ")

	if cmdContent == "clear" {
		err = bdt.Delete()
		completedEmojis = []string{"🗑", "✅"}
	} else {
		bdt.RawBioData = cmdContent

		if !hasBio {
			err = bdt.CreateRaw()
		} else {
			err = bdt.SaveRaw()
		}
	}

	if err != nil {
		return err
	}

	// react to message with a check mark to signify it worked
	for _, emoji := range completedEmojis {
		err = b.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, emoji)
		if err != nil {
			return err
		}
	}

	return nil
}
