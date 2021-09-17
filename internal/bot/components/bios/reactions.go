package bios

import "github.com/codemicro/dgo-toolkit/route"

func (b *Bios) PaginationReaction(ctx *route.ReactionContext) error {

	b.trackerLock.RLock()
	_, found := b.trackedEmbeds[ctx.Reaction.MessageID]
	b.trackerLock.RUnlock()

	if !found {
		return nil
	}

	b.trackerLock.Lock()
	defer b.trackerLock.Unlock()

	tracked := b.trackedEmbeds[ctx.Reaction.MessageID]

	if tracked.requestingUser != ctx.Reaction.UserID {
		_ = ctx.Session.MessageReactionRemove(ctx.Reaction.ChannelID, ctx.Reaction.MessageID, ctx.Reaction.Emoji.Name,
			ctx.Reaction.UserID)
		return nil
	}

	var newBioIndex int

	if ctx.Reaction.Emoji.Name == nextBioReaction {
		newBioIndex = tracked.current + 1
	} else if ctx.Reaction.Emoji.Name == previousBioReaction {
		newBioIndex = tracked.current - 1
	} else { // not one of the control emojis
		return nil
	}

	if newBioIndex < 0 || newBioIndex >= len(tracked.bios) { // out of bounds
		return nil
	}

	tracked.current = newBioIndex

	// form new embed

	plurality := &pluralityInfo{
		CurrentNumber: tracked.current + 1,
		TotalCount:    len(tracked.bios),
	}

	var nd nameDriver
	if tracked.bios[newBioIndex].SysMemberID != "" { // account bios will have a blank system member ID
		nd = newSystemName(tracked.bios[newBioIndex].SysMemberID, plurality, tracked.accountId)
	} else {
		nd = newAccountName(tracked.accountId, ctx.Reaction.GuildID, plurality, ctx.Session)
	}

	newEmbed, err := b.formBioEmbed(nd, tracked.bios[newBioIndex], tracked.isAdmin)
	if err != nil {
		return err
	}

	// remove reaction
	_ = ctx.Session.MessageReactionRemove(ctx.Reaction.ChannelID, ctx.Reaction.MessageID, ctx.Reaction.Emoji.Name, ctx.Reaction.UserID)

	// edit message with that
	_, err = ctx.Session.ChannelMessageEditEmbed(tracked.channelId, ctx.Reaction.MessageID, newEmbed)
	return err
}
