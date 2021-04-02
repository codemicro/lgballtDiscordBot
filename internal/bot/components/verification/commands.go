package verification

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"strings"
	"time"
)

func (*Verification) coreVerification(ctx *route.MessageContext) error {

	command := strings.Fields(ctx.Raw)[1:]

	// Copy message into output channel
	iu := inlineData{
		UserID: ctx.Message.Author.ID,
	}

	verificationText := strings.Join(command, " ")

	if len(verificationText) > 1500 {
		return ctx.SendErrorMessage("Sorry, that message is too long! Please keep your verification text to" +
			" a *maximum* of 1500 characters.")
	}

	var warning []string

	// check for failed verifications and bans/kicks
	var removal db.UserRemove
	var failure db.VerificationFail
	removal.UserId = ctx.Message.Author.ID
	failure.UserId = ctx.Message.Author.ID

	if found, err := removal.Get(); err != nil {
		return err
	} else if found {
		rsn := removal.Reason
		if rsn == "" {
			rsn = "none provided"
		}
		warning = append(warning, fmt.Sprintf("⚠️ **Warning**: this user has been **%s** before for reason: "+
			"*%s*\n", removal.Action, rsn))
	}

	if found, err := failure.Get(); err != nil {
		return err
	} else if found {
		warning = append(warning, fmt.Sprintf("⚠️ **Warning**: this user has failed verification before. "+
			"See %s\n", failure.MessageLink))
	}

	warningString := strings.Trim(strings.Join(warning, "\n"), "\n")

	if warningString != "" {
		warningString = "\n\n" + warningString
	}

	// form messages
	messagePartOne := fmt.Sprintf("From: %s (%s#%s)\nContent: %s", tools.MakePing(ctx.Message.Author.ID),
		ctx.Message.Author.Username, ctx.Message.Author.Discriminator, verificationText)

	var messagePartTwo string

	messagePartTwo += warningString
	messagePartTwo += "\n\n" + logHelpText
	messagePartTwo += "\n\n" + fmt.Sprintf("```%s```", iu.toString())

	var newMessage *discordgo.Message

	if len(messagePartOne)+len(messagePartTwo) > 2000 {
		_, err := ctx.SendMessageString(config.VerificationIDs.OutputChannel, messagePartOne)
		if err != nil {
			return err
		}
		newMessage, err = ctx.SendMessageString(config.VerificationIDs.OutputChannel, messagePartOne)
		if err != nil {
			return err
		}
	} else {
		var err error
		newMessage, err = ctx.SendMessageString(config.VerificationIDs.OutputChannel, messagePartOne+messagePartTwo)
		if err != nil {
			return err
		}
	}

	// Add sample reactions to that message
	for _, reaction := range []string{acceptReaction, rejectReaction} {
		err := ctx.Session.MessageReactionAdd(config.VerificationIDs.OutputChannel, newMessage.ID, reaction)
		if err != nil {
			return err
		}
	}

	// Delete user's message
	err := ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	if err != nil {
		logging.Error(err, "failed to delete message in verification")
	}

	return nil

}

func (v *Verification) Verify(ctx *route.MessageContext) error {

	command := strings.Fields(ctx.Raw)[1:]

	// check ratelimit
	if val, found := v.ratelimit[ctx.Message.Author.ID]; found && time.Now().Before(val) {
		return ctx.SendErrorMessage( "You've already submitted a verification request. Please wait.")
	}

	if len(command) < 1 {
		return ctx.SendErrorMessage("You're missing your verification message! Try again," +
			" silly! " + tools.MakeCustomEmoji(false, "trans_happy", "747448537398116392"))
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

	_, channelId, messageId, valid := tools.ParseMessageLink(messageLink)

	if !valid {
		return ctx.SendErrorMessage("Invalid message link")
	}

	mct, err := ctx.Session.ChannelMessage(channelId, messageId)
	if err != nil {
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

	err = v.coreVerification(&route.MessageContext{
		CommonContext: ctx.CommonContext,
		Message:       &discordgo.MessageCreate{Message: mct},
		Arguments:     nil,
		Raw:           mct.Content,
	})

	if err != nil {
		return err
	}

	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
}

func (*Verification) coreRecordRemoval(ctx *route.MessageContext, actionType string) error {

	userId := ctx.Arguments["user"].(string)
	reason := ctx.Arguments["reason"].(string)

	ur := db.UserRemove{UserId: userId}

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