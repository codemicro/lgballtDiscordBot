package verification

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"strings"
	"time"
)

func (v *Verification) Verify(command []string, m *harmony.Message, checkRatelimit bool) error {

	// check ratelimit
	if val, found := verificationRatelimit[m.Author.ID]; found && time.Now().Before(val) && checkRatelimit {
		_, err := v.b.SendMessage(m.ChannelID, "You've already submitted a verification request. Please wait.")
		return err
	}

	if len(command) < 1 {
		_, err := v.b.SendMessage(m.ChannelID, "You're missing your verification message! Try again, silly! "+
			tools.MakeCustomEmoji(false, "trans_happy", "747448537398116392"))
		return err
	}

	// Copy message into output channel
	iu := inlineData{
		UserID: m.Author.ID,
	}

	verificationText := strings.Join(command, " ")

	if len(verificationText) > 1500 {
		_, err := v.b.SendMessage(m.ChannelID, "Sorry, that message is too long! Please keep your "+
			"verification text to a *maximum* of 1500 characters.")
		return err
	}

	var warning string

	// check for failed verifications and bans/kicks
	var removal db.UserRemove
	var failure db.VerificationFail
	removal.UserId = m.Author.ID
	failure.UserId = m.Author.ID

	if found, err := removal.Get(); err != nil {
		return err
	} else if found {
		rsn := removal.Reason
		if rsn == "" {
			rsn = "none provided"
		}
		warning += fmt.Sprintf("⚠️ **Warning**: this user has been **%s** before for reason: *%s*\n", removal.Action,
			rsn)
	}

	if found, err := failure.Get(); err != nil {
		return err
	} else if found {
		warning += fmt.Sprintf("⚠️ **Warning**: this user has failed verification before. See %s\n", failure.MessageLink)
	}

	if warning != "" {
		warning = "\n\n" + warning
	}

	messagePartOne := fmt.Sprintf("From: %s (%s#%s)\nContent: %s", tools.MakePing(m.Author.ID), m.Author.Username, m.Author.Discriminator, verificationText)
	messagePartTwo := fmt.Sprintf("%s\n%s\n\n```%s```", warning, logHelpText, iu.toString())

	var newMessage *harmony.Message

	if len(messagePartOne)+len(messagePartTwo) > 2000 {
		_, err := v.b.SendMessage(OutputChannelId, messagePartOne)
		if err != nil {
			return err
		}
		newMessage, err = v.b.SendMessage(OutputChannelId, messagePartOne)
		if err != nil {
			return err
		}
	} else {
		var err error
		newMessage, err = v.b.SendMessage(OutputChannelId, messagePartOne+messagePartTwo)
		if err != nil {
			return err
		}
	}

	// track ratelimit
	if checkRatelimit {
		verificationRatelimit[m.Author.ID] = time.Now().Add(ratelimitTimeout)
	}

	// Add sample reactions to that message
	for _, reaction := range []string{acceptReaction, rejectReaction} {
		err := v.b.Client.Channel(OutputChannelId).AddReaction(context.Background(), newMessage.ID, reaction)
		if err != nil {
			return err
		}
	}

	// Delete user's message
	err := v.b.Client.Channel(m.ChannelID).DeleteMessage(context.Background(), m.ID)
	if err != nil {
		logging.Error(err, "failed to delete message in verification")
	}

	// Send confirmation message to user
	_, err = v.b.SendMessage(m.ChannelID, fmt.Sprintf("Thanks %s - your verification request has been "+
		"recieved. We'll check it as soon as possible.", tools.MakePing(m.Author.ID)))
	return err
}

func (v *Verification) FVerify(command []string, m *harmony.Message) error {
	if len(command) != 1 {
		_, err := v.b.SendMessage(m.ChannelID, "No message link supplied")
		return err
	}

	_, channelId, messageId, valid := tools.ParseMessageLink(command[0])

	if !valid {
		_, err := v.b.SendMessage(m.ChannelID, "Invalid message link")
		return err
	}

	mct, err := v.b.Client.Channel(channelId).Message(context.Background(), messageId)
	if err != nil {
		switch e := err.(type) {
		case *harmony.APIError:
			if e.HTTPCode == 404 {
				_, err := v.b.SendMessage(m.ChannelID, "Message not found")
				return err
			}
			return err
		default:
			return err
		}
	}

	return v.Verify(strings.Split(mct.Content, " "), mct, false)

}

func (v *Verification) RecordRemoval(command []string, m *harmony.Message) error {

	if len(command) < 3 {
		return nil
	}

	actionType := command[0]
	uid, _ := tools.ParsePing(command[1])
	reason := strings.Join(command[2:], " ")

	ur := db.UserRemove{UserId: uid}

	found, err := ur.Get()
	if err != nil {
		return err
	}

	if actionType == "ban" {
		ur.Action = "banned"
	} else {
		ur.Action = "kicked"
	}

	ur.Reason = reason

	if found {
		err = ur.Save()
	} else {
		err = ur.Create()
	}

	if err != nil {
		return err
	}

	_, err = v.b.SendMessage(m.ChannelID, "Action logged.")
	return err
}
