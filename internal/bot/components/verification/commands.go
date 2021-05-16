package verification

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/pronouns"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"strings"
	"time"
)

func (*Verification) coreVerification(ctx *route.MessageContext) error {

	command := strings.Fields(ctx.Raw)[1:]

	verificationText := strings.Join(command, " ")

	if len(verificationText) > 1500 {
		return ctx.SendErrorMessage("Sorry, that message is too long! Please keep your verification text to" +
			" a *maximum* of 1500 characters.")
	}

	guildRoles, err := ctx.Session.GuildRoles(ctx.Message.GuildID)
	if err != nil {
		return err
	}
	possiblePronouns := pronouns.FilterRoleList(guildRoles)

	emb := new(discordgo.MessageEmbed)

	// check for failed verifications and bans/kicks
	var removal db.UserRemove
	var failure db.VerificationFail
	hashedID := hashString(ctx.Message.Author.ID)
	removal.UserId = hashedID
	failure.UserId = hashedID

	if found, err := removal.Get(); err != nil {
		return err
	} else if found {
		rsn := removal.Reason
		if rsn == "" {
			rsn = "none provided"
		}

		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:   "⚠️ Warning",
			Value:  fmt.Sprintf("This user has been **%s** before for reason: *%s*", removal.Action, rsn),
			Inline: false,
		})
	}

	if found, err := failure.Get(); err != nil {
		return err
	} else if found {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:   "⚠️ Warning",
			Value:  fmt.Sprintf("This user has failed verification before. See %s", failure.MessageLink),
			Inline: false,
		})
	}

	reactionsToAdd := []string{acceptReaction, rejectReaction}

	var pronounInstructions string
	if userPronouns := pronouns.FindPronounsInString(ctx.Message.Content, possiblePronouns); len(userPronouns) != 0 {

		{
			// filter out excluded roles
			var n int
			for _, x := range userPronouns {
				if !tools.IsStringInSlice(x.RoleID, config.VerificationIDs.ExcludedPronounRoles) {
					userPronouns[n] = x
					n += 1
				}
			}
			userPronouns = userPronouns[:n]
		}

		var ps []string
		for _, p := range userPronouns {
			ps = append(ps, tools.MakeRolePing(p.RoleID))
		}

		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  pronounEmbedFieldTitle,
			Value: fmt.Sprintf("The following pronoun roles will be applied if this user is verified:\n%s", strings.Join(ps, " ")),
		})

		pronounInstructions = fmt.Sprintf("\nIf the auto-detected pronouns are incorrect, react with %s to remove them.", scrapPronounReaction)
		reactionsToAdd = append(reactionsToAdd, scrapPronounReaction)

	}

	emb.Title = fmt.Sprintf("Verification request from %s#%s", ctx.Message.Author.Username, ctx.Message.Author.Discriminator)
	emb.Description = verificationText
	emb.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("React with %s to accept this request or %s to reject this request.%s\nRejecting this request will not inform the user.", acceptReaction, rejectReaction, pronounInstructions),
	}

	newMessage, err := ctx.Session.ChannelMessageSendComplex(config.VerificationIDs.OutputChannel, &discordgo.MessageSend{
		Content: tools.MakePing(ctx.Message.Author.ID),
		Embed:   emb,
	})
	if err != nil {
		return err
	}

	// Add sample reactions to that message
	for _, reaction := range reactionsToAdd {
		err := ctx.Session.MessageReactionAdd(config.VerificationIDs.OutputChannel, newMessage.ID, reaction)
		if err != nil {
			return err
		}
	}

	// Delete user's message
	err = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	if err != nil {
		logging.Error(err, "failed to delete message in verification")
	}

	return nil

}

func (v *Verification) Verify(ctx *route.MessageContext) error {

	command := strings.Fields(ctx.Raw)[1:]

	// check ratelimit
	if val, found := v.ratelimit[ctx.Message.Author.ID]; found && time.Now().Before(val) {
		return ctx.SendErrorMessage("You've already submitted a verification request. Please wait.")
	}

	if len(command) < 1 {
		return ctx.SendErrorMessage("You're missing your verification message, which probably means you " +
			"didn't read the rules. Read " + tools.MakeChannelMention("702328069309857852") + ".")
	}

	err := v.coreVerification(ctx)
	if err != nil {
		return err
	}

	v.ratelimit[ctx.Message.Author.ID] = time.Now().Add(ratelimitTimeout)

	// Send confirmation message to user
	_, err = ctx.SendMessageString(ctx.Message.ChannelID, fmt.Sprintf("Thanks %s - your verification request has been "+
		"recieved. We'll check it as soon as possible.", tools.MakePing(ctx.Message.Author.ID)))
	return err
}

func (v *Verification) FVerify(ctx *route.MessageContext) error {

	messageLink := ctx.Arguments["messageLink"].(string)

	guildId, channelId, messageId, valid := tools.ParseMessageLink(messageLink)

	fmt.Println(channelId, messageId)

	if !valid {
		return ctx.SendErrorMessage("Invalid message link")
	}

	mct, err := ctx.Session.ChannelMessage(channelId, messageId)
	if err != nil {
		fmt.Println(err)
		switch e := err.(type) {
		case *discordgo.RESTError:
			if e.Response.StatusCode == 404 {
				return ctx.SendErrorMessage("Message not found")
			}
			return err
		default:
			return err
		}
	}

	mct.GuildID = guildId // v.coreVerification needs this, for some reason it's not set when retrieving the message
	// using ctx.Session.ChannelMessage

	err = v.coreVerification(&route.MessageContext{
		CommonContext: ctx.CommonContext,
		Message:       &discordgo.MessageCreate{Message: mct},
		Arguments:     nil,
		Raw:           mct.Content,
	})

	if err != nil {
		fmt.Println(err, "A")
		return err
	}

	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
}

func (*Verification) coreRecordRemoval(ctx *route.MessageContext, actionType string) error {

	userId := ctx.Arguments["user"].(string)
	reason := ctx.Arguments["reason"].(string)

	ur := db.UserRemove{UserId: hashString(userId)}

	found, err := ur.Get()
	if err != nil {
		return err
	}

	ur.Reason = reason
	ur.Action = actionType

	if found {
		err = ur.Save()
	} else {
		err = ur.Create()
	}

	if err != nil {
		return err
	}

	_, err = ctx.SendMessageString(ctx.Message.ChannelID, "Action logged.")
	return err
}

func (v *Verification) TrackBan(ctx *route.MessageContext) error {
	return v.coreRecordRemoval(ctx, "banned")
}

func (v *Verification) TrackKick(ctx *route.MessageContext) error {
	return v.coreRecordRemoval(ctx, "kicked")
}

func (v *Verification) PurgeUnverifiedMembers(ctx *route.MessageContext) error {

	const (
		removalThreshold = time.Hour * 24 * 7
	)

	var (
		targetGuild = ctx.Message.GuildID
		toKick      []*discordgo.Member

		lastMemberID string
		lastLength   = 1000

		kickLog strings.Builder
	)

	for lastLength >= 1000 {

		_ = ctx.Session.ChannelTyping(ctx.Message.ChannelID)

		guildMembers, err := ctx.Session.GuildMembers(targetGuild, lastMemberID, 1000)
		if err != nil {
			return err
		}

		lastLength = len(guildMembers)
		lastMemberID = guildMembers[len(guildMembers)-1].User.ID

		for _, member := range guildMembers {
			parsedJoinTime, err := member.JoinedAt.Parse()
			if err != nil {
				continue
			}

			var hasRequiredRoles bool

			for _, roleID := range member.Roles {
				if roleID == config.VerificationIDs.RoleId ||
					roleID == config.MuteMe.TimeoutRole ||
					tools.IsStringInSlice(roleID, config.VerificationIDs.ExtraValidRoles) {

					hasRequiredRoles = true
					break
				}
			}

			if !hasRequiredRoles && time.Since(parsedJoinTime) > removalThreshold && !member.User.Bot {
				toKick = append(toKick, member)
				kickLog.WriteString(member.User.ID)
				kickLog.WriteRune(' ')
				kickLog.WriteString(member.User.String())
				kickLog.WriteRune('\n')
			}
		}
	}

	emb := &discordgo.MessageEmbed{
		Title: "Confirmation",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "React to this message with ✅ to continue. (❌ to cancel)",
		},
		Description: "The following users will be kicked: ```" + kickLog.String() + "```",
	}

	kickLog.Reset()

	err := ctx.Kit.NewConfirmation(
		ctx.Message.ChannelID,
		ctx.Message.Author.ID,
		emb,
		func(ctx2 *route.ReactionContext) error {
			for i, member := range toKick {

				if i%10 == 0 {
					_ = ctx.Session.ChannelTyping(ctx.Message.ChannelID)
				}

				if !config.DebugMode {
					err := ctx.Session.GuildMemberDeleteWithReason(ctx.Message.GuildID, member.User.ID, "Automatic purge, had not verified within 7 days")
					if err != nil {
						kickLog.WriteString(member.User.ID)
						kickLog.WriteRune(' ')
						kickLog.WriteString(member.User.String())
						kickLog.WriteString(" failed: ")
						kickLog.WriteString(err.Error())
						kickLog.WriteRune('\n')
						logging.Warn(err.Error())
					}
				} else {
					logging.Info(fmt.Sprintf("KICK %s %s (debug mode enabled, no action performed)", member.User.ID, member.User.String()))
				}
			}

			var errors string
			if len(kickLog.String()) != 0 {
				errors = fmt.Sprintf(" The following errors were observed:\n```%s```", kickLog.String())
			}

			_, err := ctx.SendMessageString(ctx.Message.ChannelID, "Action(s) successful." + errors)
			return err
		},
		func(ctx *route.ReactionContext) error {
			_, _ = ctx.Session.ChannelMessageEdit(ctx.Reaction.ChannelID, ctx.Reaction.MessageID, "Cancelled!")
			return nil
		},
	)

	return err
}
