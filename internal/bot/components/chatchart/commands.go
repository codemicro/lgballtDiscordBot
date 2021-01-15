package chatchart

import (
	"context"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
)

func (c *ChatChart) TriggerCollection(command []string, m *harmony.Message) error {
	// Syntax: <channel mention>

	channelId, validChannel := tools.ParseChannelMention(command[0])
	if validChannel {
		_, err := c.b.Client.Channel(channelId).Get(context.Background())
		if err != nil {
			validChannel = false
		}
	}

	if !validChannel {
		_, err := c.b.SendMessage(m.ChannelID, "Invalid channel")
		if err != nil {
			return err
		}
	}

	var msg string

	if tools.IsStringInSlice(channelId, config.ChatChartChannelExclusions) {
		msg = "This channel has been excluded from chat chart indexing."
	} else {
		c.queue <- collectionIntent{
			ChannelId: channelId,
			Message:   m,
		}
		msg = "Task queued. You'll be pinged when collection is complete and a chart is ready."
	}

	_, err := c.b.SendMessage(m.ChannelID, msg)
	return err

}
