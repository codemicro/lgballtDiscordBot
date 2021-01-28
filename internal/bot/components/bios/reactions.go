package bios

import (
	"context"
	"github.com/skwair/harmony"
)

func (b *Bios) ReactionAdd(e *harmony.MessageReaction) error {

	b.trackerLock.RLock()
	_, found := b.trackedEmbeds[e.MessageID]
	b.trackerLock.RUnlock()

	if !found {
		return nil
	}

	b.trackerLock.Lock()
	defer b.trackerLock.Unlock()

	tracked := b.trackedEmbeds[e.MessageID]

	var newBioIndex int

	if e.Emoji.Name == nextBioReaction {
		newBioIndex = tracked.current + 1
	} else if e.Emoji.Name == previousBioReaction {
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
		nd = newSystemName(tracked.bios[newBioIndex].SysMemberID, plurality)
	} else {
		nd = newAccountName(tracked.accountId, e.GuildID, plurality, b.b)
	}

	newEmbed, err := b.formBioEmbed(nd, tracked.bios[newBioIndex].BioData)
	if err != nil {
		return err
	}

	// remove reaction
	_ = b.b.Client.Channel(tracked.channelId).RemoveUserReaction(context.Background(), e.MessageID, e.UserID, e.Emoji.Name)

	// edit message with that
	_, err = b.b.Client.Channel(tracked.channelId).EditEmbed(context.Background(), e.MessageID, "", newEmbed)
	return err
}
