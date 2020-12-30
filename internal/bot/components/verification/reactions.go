package verification

import (
	"context"
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"regexp"
	"strings"
	"time"
)

func (v *Verification) AdminDecision(r *harmony.MessageReaction) error {

	reactingUser, err := v.b.Client.User(context.Background(), r.UserID)
	if err != nil {
		return err
	}

	c := v.b.Client.Channel(r.ChannelID)

	m, err := c.Message(context.Background(), r.MessageID)
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

	if r.Emoji.Name == acceptReaction {
		actionTaken = "accepted"
		actionEmoji = acceptReaction

		err = v.b.Client.Guild(r.GuildID).AddMemberRoleWithReason(context.Background(), inlineUserData.UserID, roleId,
			fmt.Sprintf("Added verification role on request of %s#%s", reactingUser.Username,
				reactingUser.Discriminator))

		if err != nil {
			return err
		}

	} else if r.Emoji.Name == rejectReaction {
		actionTaken = "rejected"
		actionEmoji = rejectReaction

		// add verification failure
		var vf db.VerificationFail
		vf.UserId = inlineUserData.UserID
		found, err := vf.Get()
		if err != nil {
			return err
		}
		vf.MessageLink = tools.MakeMessageLink(r.GuildID, r.ChannelID, r.MessageID)
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
	newContent = strings.ReplaceAll(newContent, logHelpText, "")

	_, err = c.EditMessage(context.Background(), m.ID, newContent)
	if err != nil {
		return err
	}

	err = c.RemoveAllReactions(context.Background(), m.ID)
	return err
}
