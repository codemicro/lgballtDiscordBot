package messagetools

import (
	"context"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/hashicorp/go-multierror"
	"github.com/skwair/harmony"
	"strings"
)

func (t *MessageTools) Send(args []string, m *harmony.Message) error {
	_, err := t.b.SendMessage(args[0], strings.Join(args[1:], " "))
	emoji := "✅"
	if err != nil {
		emoji = "❌"
	}

	reactionErr := t.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, emoji)
	if reactionErr != nil {
		err = multierror.Append(err, reactionErr)
	}

	return err
}

func (t *MessageTools) Edit(args []string, m *harmony.Message) error {

	_, channelId, messageId, valid := tools.ParseMessageLink(args[0])

	if !valid {
		_, err := t.b.SendMessage(m.ChannelID, "Could not parse message link")
		return err
	}

	_, err := t.b.Client.Channel(channelId).EditMessage(context.Background(), messageId,
		strings.Join(args[1:], " "))

	emoji := "✅"
	if err != nil {
		emoji = "❌"
	}

	reactionErr := t.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, emoji)
	if reactionErr != nil {
		err = multierror.Append(err, reactionErr)
	}

	return err
}
