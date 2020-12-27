package roles

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/skwair/harmony"
)

func (r *Roles) ReactionAdd(e *harmony.MessageReaction) error {

	// check if the message has any tracked emojis
	rrs, err := db.GetAllReactionRolesForMessage(e.MessageID)
	if err != nil {
		return err
	}
	if len(rrs) == 0 {
		return nil
	}

	var emoji string
	if e.Emoji.ID != "" {
		emoji = fmt.Sprintf("%s:%s", e.Emoji.Name, e.Emoji.ID)
	} else {
		emoji = e.Emoji.Name
	}

	// check if this emoji in that list
	var emojiTracked bool
	var roleId string
	{
		for _, rr := range rrs {
			if rr.Emoji == emoji {
				emojiTracked = true
				roleId = rr.RoleId
				break
			}
		}
	}
	if !emojiTracked {
		return nil
	}

	// add role
	err = r.b.Client.Guild(e.GuildID).AddMemberRoleWithReason(context.Background(), e.UserID, roleId, "Reaction role")
	return err
}

func (r *Roles) ReactionRemove(e *harmony.MessageReaction) error {

	// check if the message has any tracked emojis
	rrs, err := db.GetAllReactionRolesForMessage(e.MessageID)
	if err != nil {
		return err
	}
	if len(rrs) == 0 {
		return nil
	}

	var emoji string
	if e.Emoji.ID != "" {
		emoji = fmt.Sprintf("%s:%s", e.Emoji.Name, e.Emoji.ID)
	} else {
		emoji = e.Emoji.Name
	}

	// check if this emoji in that list
	var emojiTracked bool
	var roleId string
	{
		for _, rr := range rrs {
			if rr.Emoji == emoji {
				emojiTracked = true
				roleId = rr.RoleId
				break
			}
		}
	}
	if !emojiTracked {
		return nil
	}

	// remove role
	err = r.b.Client.Guild(e.GuildID).RemoveMemberRole(context.Background(), e.UserID, roleId)
	return err
}
