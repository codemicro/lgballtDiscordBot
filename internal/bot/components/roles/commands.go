package roles

import (
	"context"
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/skwair/harmony"
	"regexp"
	"strings"
)

var messageLinkRegexp = regexp.MustCompile(`(?m)https://(.+)?discord\.com/channels/(\d+)/(\d+)/(\d+)/?`)

func (r *Roles) TrackReaction(command []string, m *harmony.Message) error {
	// Syntax: <message link> <emoji> <role name>

	// This is a message link: https://discord.com/channels/<guild ID>/<channel ID>/<message ID>

	// Get channel ID and message ID from link
	matches := messageLinkRegexp.FindAllStringSubmatch(command[0], -1)
	if len(matches) == 0 {
		_, err := r.b.SendMessage(m.ChannelID, "Unable to parse message link")
		if err != nil {
			return err
		}
	}

	guildId := matches[0][2]
	channelID := matches[0][3]
	messageID := matches[0][4]

	// Check message exists
	channelRes := r.b.Client.Channel(channelID)
	_, err := channelRes.Message(context.Background(), messageID)
	if err != nil {
		if errors.Is(err, harmony.APIError{}) {
			if err.(harmony.APIError).HTTPCode == 404 {
				_, err := r.b.SendMessage(m.ChannelID, "Message link invalid")
				return err
			}
		}
		return err
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

	// Search for and file role in guild
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

	// Check this role or emoji are not in use
	var errMessage string
	for _, rr := range reactionRoles {
		if rr.Emoji == command[1] {
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
	rr.Emoji = command[1]
	rr.RoleId = roleID

	// Save
	err = rr.Save()
	if err != nil {
		return err
	}

	// Confirmation
	err = r.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
	return err
}
