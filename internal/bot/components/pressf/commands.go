package pressf

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"time"
)

type activePressF struct {
	Trigger      chan *discordgo.MessageReaction
	Count        int
	ReactedUsers []string
	Message      *discordgo.Message
}


func (pf *PressF) newTracker(session *discordgo.Session, message *discordgo.Message, duration time.Duration) *activePressF {
	apf := &activePressF{
		Trigger: make(chan *discordgo.MessageReaction, 100),
		Count:   0,
		Message: message,
	}

	go func() {
		time.Sleep(duration)
		close(apf.Trigger)

		pf.mux.Lock()
		delete(pf.active, message.ID)
		pf.mux.Unlock()
	}()

	go func() {
		for v := range apf.Trigger {

			if !tools.IsStringInSlice(v.UserID, apf.ReactedUsers) {
				apf.ReactedUsers = append(apf.ReactedUsers, v.UserID)
			} else {
				continue
			}

			name, _, err := tools.GetNickname(session, v.UserID, v.GuildID)
			if err != nil {
				logging.Error(err, "activePressF runner (nickname get)")
				return
			}
			newText := fmt.Sprintf("%s\n**%s** has paid respects", message.Content, name)
			_, err = session.ChannelMessageEdit(apf.Message.ChannelID, apf.Message.ID, newText)
			apf.Message.Content = newText
			if err != nil {
				logging.Error(err, "activePressF runner")
				return
			}
			apf.Count += 1
		}
		_, err := session.ChannelMessageEdit(apf.Message.ChannelID, apf.Message.ID, fmt.Sprintf("%s\n**%d** people have paid respects", message.Content, apf.Count))
		if err != nil {
			logging.Error(err, "activePressF runner final edit")
			return
		}
		err = session.MessageReactionsRemoveAll(apf.Message.ChannelID, apf.Message.ID)
		if err != nil {
			logging.Error(err, "activePressF runner clear all reactions")
			return
		}
	}()

	return apf
}

func (pf *PressF) Trigger(ctx *route.MessageContext) error {

	payingRespectsTo := ctx.Arguments["thing"].(string)

	msg, err := ctx.SendMessageString(ctx.Message.ChannelID, fmt.Sprintf("React with 🇫 to pay your respects " +
		"to **%s**", tools.FilterRolePing(payingRespectsTo)))
	if err != nil {
		return err
	}

	err = ctx.Session.MessageReactionAdd(msg.ChannelID, msg.ID, "🇫")
	if err != nil {
		return err
	}

	apf := pf.newTracker(ctx.Session, msg, time.Minute)

	pf.mux.Lock()
	pf.active[msg.ID] = apf
	pf.mux.Unlock()

	return nil
}
