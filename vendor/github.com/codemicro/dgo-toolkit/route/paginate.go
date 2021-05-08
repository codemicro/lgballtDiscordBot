package route

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

const (
	nextReaction = "➡"
	prevReaction = "⬅"
)

func (b *Kit) NewPaginate(channelId string, userId string, embeds []*discordgo.MessageEmbed, timeout time.Duration) error {

	var current int

	msg, err := b.Session.ChannelMessageSendEmbed(channelId, embeds[current])
	if err != nil {
		return err
	}

	for _, e := range []string{prevReaction, nextReaction} {
		err = b.Session.MessageReactionAdd(msg.ChannelID, msg.ID, e)
		if err != nil {
			return err
		}
	}

	handlerId := b.AddTemporaryReaction(&Reaction{
		Name:  "Temporary pagination reaction in " + channelId,
		Run: func(ctx *ReactionContext) error {

			if ctx.Reaction.MessageID != msg.ID {
				return nil
			}

			if userId == "" || ctx.Reaction.UserID == userId {

				var delta int
				if ctx.Reaction.Emoji.Name == nextReaction {
					delta = 1
				} else if ctx.Reaction.Emoji.Name == prevReaction {
					delta = -1
				}

				if delta == 0 {
					return nil
				}

				current += delta

				if le := len(embeds); current >= le {
					current -= le
				} else if current < 0 {
					current += le
				}

				_, err = ctx.Session.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, embeds[current])
				if err != nil {
					return err
				}

				return ctx.Session.MessageReactionRemove(msg.ChannelID, msg.ID, ctx.Reaction.Emoji.Name,
					ctx.Reaction.UserID)
			}

			return nil

		},
		Event: ReactionAdd,
	})

	go func() {
		time.Sleep(timeout)
		b.RemoveTemporaryReaction(handlerId)
		_ = b.Session.MessageReactionsRemoveAll(msg.ChannelID, msg.ID)
	}()

	return nil

}