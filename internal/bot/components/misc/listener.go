package misc

import (
	"context"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/skwair/harmony"
	"sync"
	"time"
)

var (
	activeListenerRequests = make(map[string]*activeListenerRequest)
	alrMux sync.RWMutex
)


type activeListenerRequest struct {
	Trigger      chan *harmony.MessageReaction
	Message      *harmony.Message
	Bot          *core.Bot
}

const (
	acceptReaction = "✅"
	rejectReaction = "❌"
)

func newActiveListenerRequest(bot *core.Bot, message *harmony.Message, expiresIn time.Duration) *activeListenerRequest {
	alr := &activeListenerRequest{
		Trigger:   nil,
		Message: message,
		Bot: bot,
	}

	go func() {
		time.Sleep(expiresIn)
		close(alr.Trigger)

		alrMux.Lock()
		delete(activeListenerRequests, alr.Message.ID)
		alrMux.Unlock()
	}()

	go func() {
		for v := range alr.Trigger {

			// TODO: aaaaaa

		}
	}()

	return alr
}