package misc

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"strings"
)

func (s *Misc) Avatar(ctx *route.MessageContext) error {
	// Args: userID

	id := ctx.Arguments["user"].(string)

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

const bulletPoint = "•"

func (s *Misc) Help(ctx *route.MessageContext) error {

	// TODO: categories?? also paginate. and restriction warning with ⚠ emoji

	info := ctx.GetCommandInfo()

	emb := new(discordgo.MessageEmbed)

	for _, command := range info {

		var args string
		var argInfo string
		for _, arg := range command.Arguments {
			optional := arg.HasDefault

			if optional {
				args += "["
			} else {
				args += "<"
			}

			args += arg.Name

			if optional {
				args += "]"
			} else {
				args += ">"
			}

			args += " "

			argInfo += fmt.Sprintf(" %s **`%s`** - %s\n", bulletPoint, arg.Name, arg.Usage)
		}

		f := new(discordgo.MessageEmbedField)

		sep := " "
		if len(args) == 0 {
			sep = ""
		}

		f.Name = command.Name + fmt.Sprintf(" - **`%s%s%s`**", command.CommandText, sep, strings.TrimSpace(args))
		f.Value = fmt.Sprintf("%s\n%s", command.Description, argInfo)
		emb.Fields = append(emb.Fields, f)
	}

	_, err := ctx.SendMessageEmbed(ctx.Message.ChannelID, emb)
	return err
}

const (
	listenerText = "Listeners are *not* a substitute for real help and should not be treated as such. They are here only to listen to you vent, and they don’t have a obligation to help, only listen. Use the command `$ukmentalhealth`, `$usmentalhealth`, `$uslgbthelp`, and `$usrunaway` for more information."
)

func (s *Misc) ListenToMe(ctx *route.MessageContext) error {

	emb := &discordgo.MessageEmbed{
		Type:  "rich",
		Title: "Disclaimer",
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("React to this message with ✅ if you wish to ping for listeners. (❌ to cancel)"),
		},
		Description: listenerText,
	}

	err := ctx.Kit.NewConfirmation(
		ctx.Message.ChannelID,
		ctx.Message.Author.ID,
		emb,
		func(ctx *route.ReactionContext) error {

			_, err := ctx.Session.ChannelMessageSendComplex(ctx.Reaction.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("%s (for %s)", tools.MakeRolePing(config.Listeners.RoleId), tools.MakePing(ctx.Reaction.UserID)),
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Roles: []string{config.Listeners.RoleId},
				},
			})

			return err
		},
		nil,
	)

	return err
}
