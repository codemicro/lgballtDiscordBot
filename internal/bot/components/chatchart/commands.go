package chatchart

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
)

func (c *ChatChart) Trigger(ctx *route.MessageContext) error {

	channelId := ctx.Arguments["channel"].(string)

	if _, err := ctx.Session.Channel(channelId); err != nil {
		err = ctx.SendErrorMessage("Invalid channel")
		if err != nil {
			return err
		}
	}

	if tools.IsStringInSlice(channelId, config.ChatChartChannelExclusions) {
		return ctx.SendErrorMessage("This channel has been excluded from chat chart indexing.")
	} else {
		c.queue <- collectionIntent{
			ChannelId: channelId,
			Ctx:   ctx,
		}
		_, err := ctx.SendMessageString(ctx.Message.ChannelID, "Task queued. You'll be pinged when collection" +
			" is complete and a chart is ready.")
		return err
	}
}
