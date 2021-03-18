package verification

import (
	"errors"
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"regexp"
	"strings"
	"time"
)

func (*Verification) DecisionReaction(ctx *route.ReactionContext) error {

	reactingUser, err := ctx.Session.User(ctx.Reaction.UserID)
	if err != nil {
		return err
	}

	// c := v.b.Client.Channel(r.ChannelID)

	m, err := ctx.Session.ChannelMessage(ctx.Reaction.ChannelID, ctx.Reaction.MessageID)
	if err != nil {
		return err
	}

	// Fetch inline data in the message that was reacted to

	inlineUserData, err := dataFromString(m.Content)
	if err != nil {
		if errors.Is(err, errorMissingData) {
			// No data to be found, so we assume it's not relevant
			return nil
		}
		return err
	}

	// Depending on the reaction, we should do different things...
	var actionTaken string
	var actionEmoji string

	if ctx.Reaction.Emoji.Name == acceptReaction {
		actionTaken = "accepted"
		actionEmoji = acceptReaction

		err = ctx.Session.GuildMemberRoleAdd(ctx.Reaction.GuildID, inlineUserData.UserID, config.VerificationIDs.RoleId)

		if err != nil {
			return err
		}

	} else if ctx.Reaction.Emoji.Name == rejectReaction {
		actionTaken = "rejected"
		actionEmoji = rejectReaction

		// add verification failure
		var vf db.VerificationFail
		vf.UserId = inlineUserData.UserID
		found, err := vf.Get()
		if err != nil {
			return err
		}
		vf.MessageLink = tools.MakeMessageLink(ctx.Reaction.GuildID, ctx.Reaction.ChannelID, ctx.Reaction.MessageID)
		if found {
			err = vf.Save()
		} else {
			err = vf.Create()
		}
		if err != nil {
			return err
		}

	} else {
		return err
	}

	// Edit message and remove data area

	// This assumes that the data is wrapped in backticks
	re := regexp.MustCompile(`(?m)\x60\x60\x60.*\x60\x60\x60`)
	newContent := re.ReplaceAllString(m.Content, fmt.Sprintf("%s *Verification request was %s by %s at %s*",
		actionEmoji, actionTaken, tools.MakePing(reactingUser.ID), time.Now().Format("15:04 on 2 Jan 2006")))
	newContent = strings.ReplaceAll(newContent, "\n"+logHelpText+"\n", "")

	_, err = ctx.Session.ChannelMessageEdit(ctx.Reaction.ChannelID, m.ID, newContent)
	if err != nil {
		return err
	}

	err = ctx.Session.MessageReactionsRemoveAll(ctx.Reaction.ChannelID, m.ID)
	return err
}