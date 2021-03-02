package misc

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
)

const regionalFEmoji = "🇫"

var activePressFs = make(map[string]*activePressF)
var apfMux = new(sync.RWMutex)

type activePressF struct {
	Trigger      chan *harmony.MessageReaction
	Count        int
	ReactedUsers []string
	Message      *harmony.Message
	Bot          *core.Bot
}

func newPressFTracker(bot *core.Bot, message *harmony.Message, duration time.Duration) *activePressF {
	apf := &activePressF{
		Trigger: make(chan *harmony.MessageReaction, 100),
		Count:   0,
		Message: message,
		Bot:     bot,
	}

	go func() {
		time.Sleep(duration)
		close(apf.Trigger)

		apfMux.Lock()
		delete(activePressFs, message.ID)
		apfMux.Unlock()
	}()

	go func() {
		for v := range apf.Trigger {

			if !tools.IsStringInSlice(v.UserID, apf.ReactedUsers) {
				apf.ReactedUsers = append(apf.ReactedUsers, v.UserID)
			} else {
				continue
			}

			name, _, err := apf.Bot.GetNickname(v.UserID, v.GuildID)
			if err != nil {
				logging.Error(err, "activePressF runner (nickname get)")
				return
			}
			newText := fmt.Sprintf("%s\n**%s** has paid respects", message.Content, name)
			_, err = apf.Bot.Client.Channel(apf.Message.ChannelID).EditMessage(context.Background(), apf.Message.ID, newText)
			apf.Message.Content = newText
			if err != nil {
				logging.Error(err, "activePressF runner")
				return
			}
			apf.Count += 1
		}
		_, err := apf.Bot.Client.Channel(apf.Message.ChannelID).EditMessage(context.Background(), apf.Message.ID, fmt.Sprintf("%s\n**%d** people have paid respects", message.Content, apf.Count))
		if err != nil {
			logging.Error(err, "activePressF runner final edit")
			return
		}
		err = apf.Bot.Client.Channel(apf.Message.ChannelID).RemoveAllReactions(context.Background(), apf.Message.ID)
		if err != nil {
			logging.Error(err, "activePressF runner clear all reactions")
			return
		}
	}()

	return apf
}

func (s *Misc) PressF(command []string, m *harmony.Message) error {

	if len(command) == 0 {
		_, err := s.b.SendMessage(m.ChannelID, "You've forgotten to include what you're paying respects to. Try again!")
		return err
	}

	payingRespectsTo := strings.Join(command, " ")

	msg, err := s.b.SendMessage(m.ChannelID, fmt.Sprintf("React with 🇫 to pay your respects to **%s**",
		tools.FilterRolePing(payingRespectsTo)))
	if err != nil {
		return err
	}

	err = s.b.Client.Channel(msg.ChannelID).AddReaction(context.Background(), msg.ID, "🇫")
	if err != nil {
		return err
	}

	apf := newPressFTracker(s.b, msg, time.Minute)

	apfMux.Lock()
	activePressFs[msg.ID] = apf
	apfMux.Unlock()

	return nil
}

func (s *Misc) PressFReaction(r *harmony.MessageReaction) error {
	apfMux.RLock()
	apf, found := activePressFs[r.MessageID]
	apfMux.RUnlock()
	if found && r.Emoji.Name == regionalFEmoji {
		apf.Trigger <- r
	}
	return nil
}
