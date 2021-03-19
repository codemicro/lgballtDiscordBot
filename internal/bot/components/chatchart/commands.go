package chatchart

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
)

func (c *ChatChart) Trigger(ctx *route.MessageContext) error {

	channelId := ctx.Arguments["channel"].(string)

	if _, err := ctx.Session.Channel(channelId); err != nil {
		_, err := ctx.SendMessageString(ctx.Message.ChannelID, "Invalid channel")
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
			Ctx:   ctx,
		}
		msg = "Task queued. You'll be pinged when collection is complete and a chart is ready."
	}

	_, err := ctx.SendMessageString(ctx.Message.ChannelID, msg)
	return err

}
