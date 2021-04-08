package verification

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
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

	if m.Author.ID != ctx.Session.State.User.ID || m.ChannelID != config.VerificationIDs.OutputChannel || len(m.Embeds) < 1 {
		return nil
	}

	// Fetch inline data in the message that was reacted to

	// find the user ID that this relates to
	userID, found := tools.ParsePing(m.Content)
	if !found {
		return nil
	}

	// Depending on the reaction, we should do different things...
	var actionTaken string
	var actionEmoji string

	if ctx.Reaction.Emoji.Name == acceptReaction {
		actionTaken = "accepted"
		actionEmoji = acceptReaction

		err = ctx.Session.GuildMemberRoleAdd(ctx.Reaction.GuildID, userID, config.VerificationIDs.RoleId)

		if err != nil {
			return err
		}

	} else if ctx.Reaction.Emoji.Name == rejectReaction {
		actionTaken = "rejected"
		actionEmoji = rejectReaction

		// add verification failure
		var vf db.VerificationFail
		vf.UserId = userID
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

	// Edit message

	m.Embeds[0].Fields = append(m.Embeds[0].Fields, &discordgo.MessageEmbedField{
		Name:   fmt.Sprintf("%s Decision", actionEmoji),
		Value:  fmt.Sprintf("Verification request was %s by %s at %s", actionTaken, tools.MakePing(reactingUser.ID), time.Now().Format(time.RFC822)),
		Inline: false,
	})

	m.Embeds[0].Footer = nil

	_, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content:         &m.Content,
		Embed:           m.Embeds[0],
		ID:              m.ID,
		Channel:         m.ChannelID,
	})

	err = ctx.Session.MessageReactionsRemoveAll(ctx.Reaction.ChannelID, m.ID)
	return err
}