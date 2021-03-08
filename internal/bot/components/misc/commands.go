package misc

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
)

func (s *Misc) Avatar(ctx *route.MessageContext) error {
	// Args: userID

	id := ctx.Arguments["userId"].(string)

	// get user
	user, err := ctx.Session.User(id)
	if err != nil {
		switch n := err.(type) {
		case *discordgo.RESTError:
			if n.Response.StatusCode == 404 || n.Response.StatusCode == 400 {
				_, err := ctx.SendMessageString(ctx.Message.ChannelID, "This user doesn't exist.")
				return err
			}
		}
		return err
	}

	// send message
	_, err = ctx.SendMessageString(ctx.Message.ChannelID, user.AvatarURL(""))
	return err
}

func (s *Misc) coreEmojiCommand(ctx *route.MessageContext, emoji string) error {
	validEmoji, animated, _, emojiID := tools.ParseEmojiComponents(emoji)

	if !validEmoji {
		_, err := ctx.SendMessageString(ctx.Message.ChannelID, "That's not a valid (custom) emoji!")
		return err
	}

	emojiUrl := "https://cdn.discordapp.com/emojis/" + emojiID
	if animated {
		emojiUrl += ".gif"
	} else {
		emojiUrl += ".png"
	}

	_, err := ctx.SendMessageString(ctx.Message.ChannelID, fmt.Sprintf("ID: `%s`\n%s", emojiID, emojiUrl))
	return err
}

func (s *Misc) Emoji(ctx *route.MessageContext) error {
	// Args: emoji
	return s.coreEmojiCommand(ctx, ctx.Arguments["emoji"].(string))
}

func (s *Misc) StealEmojis(ctx *route.MessageContext) error {
	// Args: messageLink

	_, channelId, messageId, validMessageLink := tools.ParseMessageLink(ctx.Arguments["messageLink"].(string))

	if !validMessageLink {
		_, err := ctx.SendMessageString(ctx.Message.ChannelID, "Invalid message link")
		return err
	}

	msg, err := ctx.Session.ChannelMessage(channelId, messageId)
	if err != nil {
		return err
	}

	emojisFromMessage := tools.CustomEmojiRegex.FindAllString(msg.Content, -1)

	for _, emoji := range emojisFromMessage {
		err := s.coreEmojiCommand(ctx, emoji)
		if err != nil {
			return err
		}
	}

	if len(emojisFromMessage) == 0 {
		_, err = ctx.SendMessageString(ctx.Message.ChannelID, "No custom emojis found in that message")
	}

	return err
}