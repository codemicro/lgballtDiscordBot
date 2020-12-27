package core

import (
	"context"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"sync"
)

type Bot struct {
	Client *harmony.Client
	Prefix string
}

func New(client *harmony.Client, prefix string) *Bot {
	b := new(Bot)
	b.Client = client
	b.Prefix = prefix

	return b

}

func (b *Bot) SendMessage(channelID, content string) (*harmony.Message, error) {
	return b.Client.Channel(channelID).SendMessage(context.Background(), content)
}

func (b *Bot) SendEmbed(channelID string, e *embed.Embed) (*harmony.Message, error) {
	return b.Client.Channel(channelID).Send(context.Background(), harmony.WithEmbed(e))
}

var ownId string
var idOnce sync.Once

func (b *Bot) IsSelf(id string) (bool, error) {
	var err error
	idOnce.Do(func() {
		var usr *harmony.User
		usr, err = b.Client.CurrentUser().Get(context.Background())
		ownId = usr.ID
	})
	if err != nil {
		return false, err
	}
	return id == ownId, nil
}
