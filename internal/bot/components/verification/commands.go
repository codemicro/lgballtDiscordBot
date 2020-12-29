package verification

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"strings"
)

func (v *Verification) Verify(command []string, m *harmony.Message) error {
	// Copy message into output channel
	iu := inlineData{
		UserID:    m.Author.ID,
	}

	verificationText := strings.Join(command, " ")

	if len(verificationText) > 1500 {
		_, err := v.b.SendMessage(m.ChannelID, "Sorry, that message is too long! Please keep your " +
			"verification text to a *maximum* of 1500 characters.")
		return err
	}

	message := fmt.Sprintf("From: %s (%s#%s)\nContent: %s%s\n\n```%s```", tools.MakePing(m.Author.ID),
		m.Author.Username, m.Author.Discriminator, verificationText, logHelpText, iu.toString())

	newMessage, err := v.b.SendMessage(OutputChannelId, message)
	if err != nil {
		return err
	}

	// Add sample reactions to that message
	for _, reaction := range []string{acceptReaction, rejectReaction} {
		err = v.b.Client.Channel(OutputChannelId).AddReaction(context.Background(), newMessage.ID, reaction)
		if err != nil {
			return err
		}
	}

	// Delete user's message
	err = v.b.Client.Channel(m.ChannelID).DeleteMessage(context.Background(), m.ID)
	if err != nil {
		logging.Error(err, "failed to delete message in verification")
	}

	// Send confirmation message to user
	_, err = v.b.SendMessage(m.ChannelID, fmt.Sprintf("Thanks %s - your verification request has been " +
		"recieved. We'll check it as soon as possible.", tools.MakePing(m.Author.ID)))
	return err
}
