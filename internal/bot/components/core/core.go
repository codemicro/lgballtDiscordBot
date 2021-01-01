package core

import (
	"context"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"regexp"
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

var pingRegex = regexp.MustCompile(`(?m)@(?:everyone|here)`)

const pingRegexSub = "`$0`"

func (b *Bot) SendMessage(channelID, content string) (*harmony.Message, error) {
	return b.Client.Channel(channelID).SendMessage(context.Background(),
		pingRegex.ReplaceAllString(content, pingRegexSub))
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

func (b *Bot) GetNickname(uid, guildId string) (string, *harmony.User, error) {
	// Attempt to get guild member
	member, err := b.Client.Guild(guildId).Member(context.Background(), uid)
	if err != nil {
		switch e := err.(type) {
		case *harmony.APIError:
			if e.HTTPCode == 404 {
				// Can't get the member, so just get the user instead
				user, err := b.Client.User(context.Background(), uid)
				if err != nil {
					return "", nil, err
				}
				return user.Username, user, nil
			}
		default:
			return "", nil, err
		}
	} else {
		name := member.Nick
		if name == "" {
			name = member.User.Username
		}
		return name, member.User, nil
	}
	return "", nil, nil
}
