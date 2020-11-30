package bios

import (
	"context"
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"regexp"
	"strings"
)

var idFromPingRegex = regexp.MustCompile(`(?m)<@!(.+)>`)

func (b *bios) onMessageCreate(m *harmony.Message) {

	if isSelf, err := b.b.IsSelf(m.Author.ID); err != nil {
		logging.Error(err)
		return
	} else if isSelf {
		return
	}

	// ignore bots
	if m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, b.commandPrefix) {

		command := tools.GetCommand(m.Content, b.commandPrefix)

		if len(command) == 0 {
			// Get requesting user's bio
			err := performBioAction(b, m, m.Author.ID)
			if err != nil {
				logging.Error(err)
			}

		} else if len(command) == 1 {
			if strings.EqualFold(command[0], "help") {
				_, err := b.b.SendEmbed(m.ChannelID, biosHelpEmbed)
				if err != nil {
					logging.Error(err)
				}
				return
			}

			// Get bio of user with ID of first argument
			// Syntax: <uid>

			// If there's a ping as the argument, use the ID from that. Else, just use the plain argument
			id := command[0]
			if v := idFromPingRegex.FindStringSubmatch(command[0]); len(v) > 1 {
				id = v[1]
			}

			err := performBioAction(b, m, id)
			if err != nil {
				logging.Error(err)
			}
		} else if len(command) >= 2 {
			// Set requesting user's bio field to the named argument
			// Syntax: <field> <value>

			var validFieldName bool
			var properFieldName string // This is so the capitalisation is correct. In reality, I probably don't need
			// this and I could do it *properly*, but what's the point. This works.
			// properFieldName is used when setting the value of a field in the bio map.
			{
				for _, f := range b.data.Fields {
					if strings.EqualFold(f, command[0]) {
						validFieldName = true
						properFieldName = f
						break
					}
				}
			}

			if !validFieldName {
				_, err := b.b.SendMessage(m.ChannelID, "That's not a valid field name! Choose from one of the " +
					"following: " + strings.Join(b.data.Fields, ", "))

				if err != nil {
					logging.Error(err)
				}

				return
			}

			b.data.Lock.RLock()
			bioData, hasBio := b.data.UserBios[m.Author.ID]
			b.data.Lock.RUnlock()

			newValue := strings.Join(command[1:], " ")

			if !hasBio {
				bioData = make(map[string]string)
				bioData[properFieldName] = newValue
			} else {
				bioData[properFieldName] = newValue
			}

			b.data.Lock.Lock()
			b.data.UserBios[m.Author.ID] = bioData
			b.data.Lock.Unlock()

			// react to message with a check mark to signify it worked
			err := b.b.Client.Channel(m.ChannelID).AddReaction(context.Background(), m.ID, "✅")
			if err != nil {
				logging.Error(err)
				return
			}

		}

	}

}

func performBioAction(b *bios, m *harmony.Message, uid string) error {
	b.data.Lock.RLock()
	_, ok := b.data.UserBios[uid]
	b.data.Lock.RUnlock()

	if !ok {
		// No bio for that user
		_, err := b.b.SendMessage(m.ChannelID, "This user hasn't created a bio, or just plain doesn't exist.")
		if err != nil {
			return err
		}
	} else {
		// Found a bio, now to form an embed
		e, err := formBioEmbed(b, uid)
		if err != nil {
			return err
		}

		_, err = b.b.SendEmbed(m.ChannelID, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func formBioEmbed(b *bios, uid string) (*embed.Embed, error) {

	// Get user that owns this bio
	bioUser, err := b.b.Client.User(context.Background(), uid)
	if err != nil {
		return &embed.Embed{}, err
	}

	// Fetch that user's bio
	b.data.Lock.RLock()
	val, ok := b.data.UserBios[bioUser.ID]
	b.data.Lock.RUnlock()

	if !ok {
		return &embed.Embed{}, errors.New("the specified user has no bio. If you're seeing this, you've made a" +
			" programming error, you fool")
	}

	e := embed.New()
	e.Thumbnail(embed.NewThumbnail(bioUser.AvatarURL()))
	e.Title(fmt.Sprintf("%s's bio", bioUser.Username))

	var fields []*embed.Field
	for _, category := range b.data.Fields {
		fVal, ok := val[category]
		if ok {
			fields = append(fields, embed.NewField().Name(category).Value(fVal).Build())
		}
	}

	e.Fields(fields...)

	return e.Build(), nil
}
