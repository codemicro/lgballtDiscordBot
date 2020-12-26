package roles

import (
	"context"
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"strings"
)

func (r *Roles) TrackReaction(command []string, m *harmony.Message) error {
	// Syntax: <message link> <emoji> <role name>

	guildId, channelID, messageID, valid := tools.ParseMessageLink(command[0])

	if !valid {
		_, err := r.b.SendMessage(m.ChannelID, "Unable to parse message link")
		return err
	}

	// Check message exists
	channelRes := r.b.Client.Channel(channelID)
	_, err := channelRes.Message(context.Background(), messageID)
	if err != nil {
		var respCode int
		switch e := err.(type) {
		case *harmony.APIError:
			respCode = e.HTTPCode
		case *harmony.ValidationError:
			respCode = e.HTTPCode
		default:
			return err
		}

		if respCode == 404 || respCode == 400 {
			_, err := r.b.SendMessage(m.ChannelID, "Message link invalid")
			return err
		}
	}

	// Get role reactions for message
	reactionRoles, err := db.GetAllReactionRolesForMessage(messageID)
	if err != nil {
		return err
	}

	// Check there are less than 20

	if len(reactionRoles) >= 20 {
		_, err := r.b.SendMessage(m.ChannelID, "This message has too many reaction roles assigned to it " +
			"already (maximum 20)")
		return err
	}

	// Search for and find role in guild
	guild := r.b.Client.Guild(guildId)
	roles, err := guild.Roles(context.Background())
	if err != nil {
		if errors.Is(err, harmony.APIError{}) {
			_, err := r.b.SendMessage(m.ChannelID, "Unable to fetch roles for this guild")
			return err
		}
		return err
	}

	var roleID string
	targetRoleName := strings.Join(command[2:], " ")
	{
		for _, role := range roles {
			if strings.EqualFold(role.Name, targetRoleName) {
				roleID = role.ID
				break
			}
		}
	}
	if roleID == "" {
		_, err := r.b.SendMessage(m.ChannelID, "Unable to find the specified role")
		return err
	}

	// Determine emoji string
	emoji := tools.ParseEmojiToString(command[1])

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
		_, err := r.b.SendMessage(m.ChannelID, errMessage)
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

	// React on message

	err = r.b.Client.Channel(channelID).AddReaction(context.Background(), messageID, strings.TrimRight(strings.Join(strings.Split(command[1], ":")[1:3], ":"), ">"))
	if err != nil {
		return err
	}

	// Confirmation
	return r.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
}

func (r *Roles) UntrackReaction(command []string, m *harmony.Message) error {
	// Syntax: <message link> <emoji>

	_, channelID, messageID, valid := tools.ParseMessageLink(command[0])

	if !valid {
		_, err := r.b.SendMessage(m.ChannelID, "Unable to parse message link")
		return err
	}

	// Check message exists
	channelRes := r.b.Client.Channel(channelID)
	_, err := channelRes.Message(context.Background(), messageID)
	if err != nil {
		var respCode int
		switch e := err.(type) {
		case *harmony.APIError:
			respCode = e.HTTPCode
		case *harmony.ValidationError:
			respCode = e.HTTPCode
		default:
			return err
		}

		if respCode == 404 || respCode == 400 {
			_, err := r.b.SendMessage(m.ChannelID, "Message link invalid")
			return err
		}
	}

	// NUKE!
	rr := &db.ReactionRole{
		MessageId: messageID,
		Emoji:     tools.ParseEmojiToString(command[1]),
	}
	err = rr.Delete()
	if err != nil {
		return err
	}

	return r.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
}
