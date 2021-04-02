package roles

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"strings"
)

func (*Roles) Track(ctx *route.MessageContext) error {
	// Syntax: <message link> <emoji> <role name>

	messageLink := ctx.Arguments["messageLink"].(string)
	emojiString := ctx.Arguments["emoji"].(string)
	roleName := ctx.Arguments["roleName"].(string)

	guildId, channelID, messageID, valid := tools.ParseMessageLink(messageLink)

	if !valid {
		return ctx.SendErrorMessage("Unable to parse message link")
	}

	// Check message exists
	_, err := ctx.Session.ChannelMessage(channelID, messageID)
	if err != nil {
		var respCode int
		switch e := err.(type) {
		case *discordgo.RESTError:
			respCode = e.Response.StatusCode
		default:
			return err
		}

		if respCode == 404 || respCode == 400 {
			return ctx.SendErrorMessage("Message link invalid")
		}
	}

	// Get role reactions for message
	reactionRoles, err := db.GetAllReactionRolesForMessage(messageID)
	if err != nil {
		return err
	}

	// Check there are less than 20

	if len(reactionRoles) >= 20 {
		return ctx.SendErrorMessage("This message has too many reaction roles assigned to it already " +
			"(maximum 20)")
	}

	// Search for and find role in guild

	roles, err := ctx.Session.GuildRoles(guildId)
	if err != nil {
		return err
	}

	var roleID string
	{
		for _, role := range roles {
			if strings.EqualFold(role.Name, roleName) {
				roleID = role.ID
				break
			}
		}
	}
	if roleID == "" {
		return ctx.SendErrorMessage("Unable to find the specified role")
	}

	// Determine emoji string
	emoji := tools.ParseEmojiToString(emojiString)

	// Check this role or emoji are not in use
	var errMessage string
	for _, rr := range reactionRoles {
		if rr.Emoji == emoji {
			errMessage = "Emoji already in use on this message"
			break
		}
		if rr.RoleId == roleID {
			errMessage = fmt.Sprintf("Role is already assigned to %s on this message", rr.Emoji)
			break
		}
	}
	if errMessage != "" {
		return ctx.SendErrorMessage(errMessage)
	}

	// Add reaction
	// This will also act as a litmus test to see if the emoji is actually valid, since this will return an error if the
	// emoji is not valid
	err = ctx.Session.MessageReactionAdd(channelID, messageID, emoji)
	if err != nil {
		switch n := err.(type) {
		case *discordgo.RESTError:
			if n.Response.StatusCode == 400 && n.Message != nil &&
				strings.Contains(strings.ToLower(n.Message.Message), "unknown emoji") {
				return ctx.SendErrorMessage("That's not a valid emoji")
			}
		}
		return err
	}

	// Create new object
	rr := new(db.ReactionRole)
	rr.MessageId = messageID
	rr.Emoji = emoji
	rr.RoleId = roleID

	// Save
	err = rr.Create()
	if err != nil {
		return err
	}

	// Confirmation
	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
}

func (r *Roles) Untrack(ctx *route.MessageContext) error {
	// Syntax: <message link> <emoji>

	messageLink := ctx.Arguments["messageLink"].(string)
	emojiString := ctx.Arguments["emoji"].(string)

	_, channelID, messageID, valid := tools.ParseMessageLink(messageLink)

	if !valid {
		return ctx.SendErrorMessage("Unable to parse message link")
	}

	// Check message exists
	_, err := ctx.Session.ChannelMessage(channelID, messageID)
	if err != nil {
		var respCode int
		switch e := err.(type) {
		case *discordgo.RESTError:
			respCode = e.Response.StatusCode
		default:
			return err
		}

		if respCode == 404 || respCode == 400 {
			return ctx.SendErrorMessage("Message link invalid")
		}
	}

	// NUKE!
	rr := &db.ReactionRole{
		MessageId: messageID,
		Emoji:     tools.ParseEmojiToString(emojiString),
	}
	err = rr.Delete()
	if err != nil {
		return err
	}

	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
}
