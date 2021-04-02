package messageTools

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/hashicorp/go-multierror"
)

func (*MessageTools) Send(ctx *route.MessageContext) error {

	channelId := ctx.Arguments["channelId"].(string)
	message := ctx.Arguments["message"].(string)

	_, err := ctx.SendMessageString(channelId, message)
	emoji := "✅"
	if err != nil {
		emoji = "❌"
	}

	reactionErr := ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, emoji)
	if reactionErr != nil {
		err = multierror.Append(err, reactionErr)
	}

	return err
}

func (*MessageTools) Edit(ctx *route.MessageContext) error {

	messageLink := ctx.Arguments["messageLink"].(string)
	newContent := ctx.Arguments["newContent"].(string)

	_, channelId, messageId, valid := tools.ParseMessageLink(messageLink)

	if !valid {
		return ctx.SendErrorMessage("Could not parse message link")
	}

	_, err := ctx.Session.ChannelMessageEdit(channelId, messageId, newContent)

	emoji := "✅"
	if err != nil {
		emoji = "❌"
	}

	reactionErr := ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, emoji)
	if reactionErr != nil {
		err = multierror.Append(err, reactionErr)
	}

	return err
}