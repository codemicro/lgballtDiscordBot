package misc

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
)

func (s *Misc) Avatar(command []string, m *harmony.Message) error {
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

	// get user
	user, err := s.b.Client.User(context.Background(), id)
	if err != nil {
		switch n := err.(type) {
		case *harmony.APIError:
			if n.HTTPCode == 404 {
				_, err := s.b.SendMessage(m.ChannelID, "This user doesn't exist.")
				return err
			}
		}
		return err
	}

	// send message
	_, err = s.b.SendMessage(m.ChannelID, user.AvatarURL())
	return err
}

func (s *Misc) Emoji(command []string, m *harmony.Message) error {
	// Syntax: <emoji>

	validEmoji, animated, _, emojiID := tools.ParseEmojiComponents(command[0])

	if !validEmoji {
		_, err := s.b.SendMessage(m.ChannelID, "That's not a valid (custom) emoji!")
		return err
	}

	emojiUrl := "https://cdn.discordapp.com/emojis/" + emojiID
	if animated {
		emojiUrl += ".gif"
	} else {
		emojiUrl += ".png"
	}

	_, err := s.b.SendMessage(m.ChannelID, fmt.Sprintf("ID: `%s`\n%s", emojiID, emojiUrl))
	return err
}
