package roles

import (
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
)

func (r *Roles) commonReactionHandler(ctx *route.ReactionContext, isAdd bool) error {

	// check if the message has any tracked emojis
	rrs, err := db.GetAllReactionRolesForMessage(ctx.Reaction.MessageID)
	if err != nil {
		return err
	}
	if len(rrs) == 0 {
		return nil
	}

	var emoji string
	if ctx.Reaction.Emoji.ID != "" {
		emoji = fmt.Sprintf("%s:%s", ctx.Reaction.Emoji.Name, ctx.Reaction.Emoji.ID)
	} else {
		emoji = ctx.Reaction.Emoji.Name
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

	if isAdd {
		// add role
		return ctx.Session.GuildMemberRoleAdd(ctx.Reaction.GuildID, ctx.Reaction.UserID, roleId)
	}
	// remove role
	return ctx.Session.GuildMemberRoleRemove(ctx.Reaction.GuildID, ctx.Reaction.UserID, roleId)
}

func (r *Roles) ReactionAdd(ctx *route.ReactionContext) error {
	return r.commonReactionHandler(ctx, true)
}

func (r *Roles) ReactionRemove(ctx *route.ReactionContext) error {
	return r.commonReactionHandler(ctx, false)
}
