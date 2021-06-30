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

	m, err := ctx.Session.ChannelMessage(ctx.Reaction.ChannelID, ctx.Reaction.MessageID)
	if err != nil {
		return err
	}

	if m.Author.ID != ctx.Session.State.User.ID || m.ChannelID != config.VerificationIDs.OutputChannel || len(m.Embeds) < 1 {
		return nil
	}

	if ctx.Reaction.Emoji.Name == scrapPronounReaction {
		n := 0
		for _, field := range m.Embeds[0].Fields {
			if field.Name != pronounEmbedFieldTitle {
				m.Embeds[0].Fields[n] = field
				n++
			}
		}
		m.Embeds[0].Fields = m.Embeds[0].Fields[:n]

		if _, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content: &m.Content,
			Embed:   m.Embeds[0],
			ID:      m.ID,
			Channel: m.ChannelID,
		}); err != nil {
			return err
		}

		return ctx.Session.MessageReactionsRemoveEmoji(ctx.Reaction.ChannelID, ctx.Reaction.MessageID, scrapPronounReaction)

	}

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

		var pronounContent string
		for _, field := range m.Embeds[0].Fields {
			if field.Name == pronounEmbedFieldTitle {
				pronounContent = field.Value
				break
			}
		}

		rolesToAdd := []string{config.VerificationIDs.RoleId}

		if pronounContent != "" {
			rolesToAdd = append(rolesToAdd, tools.FindRolePings(pronounContent)...)
		}

		for _, role := range rolesToAdd {
			if err = ctx.Session.GuildMemberRoleAdd(ctx.Reaction.GuildID, userID, role); err != nil {
				// typically, this means that the user we're referring to has left the guild
				// in that case, let's clear the embed

				if de, ok := err.(*discordgo.RESTError); ok {
					if de.Response.StatusCode == 404 {
						// user not found (has left?)
						m.Embeds[0].Fields = append(m.Embeds[0].Fields, &discordgo.MessageEmbedField{
							Name:   "🤔 User not found",
							Value:  fmt.Sprintf("It looks like %s left before a decision was reached.", tools.MakePing(userID)),
							Inline: false,
						})

						m.Embeds[0].Footer = nil

						var x string
						_, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
							Content: &x,
							Embed:   m.Embeds[0],
							ID:      m.ID,
							Channel: m.ChannelID,
						})

						return ctx.Session.MessageReactionsRemoveAll(ctx.Reaction.ChannelID, m.ID)
					}
				}

				_, _ = ctx.SendMessageString(ctx.Reaction.ChannelID, err.Error())
				return err
			}
		}

	} else if ctx.Reaction.Emoji.Name == rejectReaction {
		actionTaken = "rejected"
		actionEmoji = rejectReaction

		// add verification failure
		var vf db.VerificationFail
		vf.UserId = hashString(userID)
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
		Value:  fmt.Sprintf("%s was %s by %s at %s", tools.MakePing(userID), actionTaken, tools.MakePing(reactingUser.ID), time.Now().Format(time.RFC822)),
		Inline: false,
	})

	m.Embeds[0].Footer = nil

	var x string
	_, err = ctx.Session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content: &x,
		Embed:   m.Embeds[0],
		ID:      m.ID,
		Channel: m.ChannelID,
	})

	return ctx.Session.MessageReactionsRemoveAll(ctx.Reaction.ChannelID, m.ID)
}
