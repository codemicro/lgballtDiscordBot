package misc

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"sync"
)

var (
	activeListenerRequests = make(map[string]struct{})
	alrMux                 sync.RWMutex
)

const (
	listenerAcceptReaction = "✅"
	listenerRejectReaction = "❌"
)

var (
	listenerText = "Listeners are *not* a substitute for real help and should not be treated as such. They are here only to listen to you vent, and they don’t have a obligation to help, only listen. Use the command `$ukmentalhealth`, `$usmentalhealth`, `$uslgbthelp`, and `$usrunaway` for more information."
)

func (s *Misc) ListenToMe(_ []string, m *harmony.Message) error {

	var found bool
	for _, cid := range config.Listeners.AllowedChannels {
		if cid == m.ChannelID {
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	emb := embed.Embed{
		Type:        "rich",
		Title:       "Disclaimer",
		Footer:      embed.NewFooter().Text(fmt.Sprintf("React to this message with %s if you wish to ping for listeners. (%s to cancel)", listenerAcceptReaction, listenerRejectReaction)).Build(),
		Description: listenerText,
	}

	msg, err := s.b.SendEmbed(m.ChannelID, &emb)
	if err != nil {
		return err
	}

	alrMux.Lock()
	activeListenerRequests[msg.ID] = struct{}{}
	alrMux.Unlock()

	for _, v := range []string{listenerAcceptReaction, listenerRejectReaction} {
		err := s.b.Client.Channel(msg.ChannelID).AddReaction(context.Background(), msg.ID, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Misc) ListenerReaction(r *harmony.MessageReaction) error {

	if !(r.Emoji.Name == listenerAcceptReaction || r.Emoji.Name == listenerRejectReaction) {
		return nil
	}

	alrMux.RLock()
	_, found := activeListenerRequests[r.MessageID]
	alrMux.RUnlock()

	if !found {
		return nil
	}

	err := s.b.Client.Channel(r.ChannelID).RemoveAllReactions(context.Background(), r.MessageID)
	if err != nil {
		return err
	}

	alrMux.Lock()
	delete(activeListenerRequests, r.MessageID)
	alrMux.Unlock()

	if r.Emoji.Name == listenerAcceptReaction {
		_, err = s.b.SendMessage(r.ChannelID, fmt.Sprintf("%s (for %s)", tools.MakeRolePing(config.Listeners.RoleId),
			tools.MakePing(r.UserID)))
		if err != nil {
			return err
		}
	}

	return nil
}
