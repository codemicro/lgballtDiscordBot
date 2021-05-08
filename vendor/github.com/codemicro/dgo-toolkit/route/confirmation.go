package route

import (
	"github.com/bwmarrin/discordgo"
)

const acceptReaction = "✅"
const rejectReaction = "❌"

// NewConfirmation will create a new confirmation embed in the channel with the ID channelId. If userID is specified,
// reactions by users other than the user with that ID are ignored - else, all reactions are accepted. acceptFunc and
// rejectFunc can also be nil.
func (b *Kit) NewConfirmation(channelId string, userId string, embed *discordgo.MessageEmbed, acceptFunc ReactionRunFunc, rejectFunc ReactionRunFunc) error {

	msg, err := b.Session.ChannelMessageSendEmbed(channelId, embed)
	if err != nil {
		return err
	}

	handlerId := -1
	handlerId = b.AddTemporaryReaction(&Reaction{
		Name: "Temporary reaction in " + channelId,
		Run: func(ctx *ReactionContext) error {
			if ctx.Reaction.MessageID != msg.ID {
				return nil
			}

			if userId == "" || ctx.Reaction.UserID == userId {

				var err error

				if ctx.Reaction.Emoji.Name == acceptReaction {
					if acceptFunc != nil {
						err = acceptFunc(ctx)
					}
				} else if ctx.Reaction.Emoji.Name == rejectReaction {
					if rejectFunc != nil {
						err = rejectFunc(ctx)
					}
				} else {
					return nil
				}

				if err != nil {
					return err
				}

				go b.RemoveTemporaryReaction(handlerId) // goroutine because it'll deadlock otherwise

				err = ctx.Session.MessageReactionsRemoveAll(channelId, msg.ID)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Event: ReactionAdd,
	})

	for _, r := range []string{acceptReaction, rejectReaction} {
		err = b.Session.MessageReactionAdd(channelId, msg.ID, r)
		if err != nil {
			return err
		}
	}

	return nil

}
