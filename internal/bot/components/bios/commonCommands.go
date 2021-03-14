package bios

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"time"
)

func (b *Bios) ReadBio(ctx *route.MessageContext) error {

	targetUserId := ctx.Arguments["user"].(string)

	bios, err := db.GetBiosForAccount(targetUserId)
	if err != nil {
		return err
	}

	if len(bios) == 0 {
		_, err := ctx.SendMessageString(ctx.Message.ChannelID, "This user hasn't created a bio, or just plain doesn't exist.")
		return err
	}

	if len(bios) == 1 {
		// Found a bio, now to form an embed
		e, err := b.formBioEmbed(newAccountName(targetUserId, ctx.Message.GuildID, nil, ctx.Session), bios[0].BioData)
		if err != nil {
			return err
		}

		_, err = ctx.SendMessageEmbed(ctx.Message.ChannelID, e)
		if err != nil {
			return err
		}
	} else {

		totalBios := len(bios)

		tracker := &trackedEmbed{
			accountId: targetUserId,
			channelId: ctx.Message.ChannelID,
			bios:      bios,
			timeoutAt: time.Now().Add(bioTimeoutDuration),
		}

		// send first bio

		plurality := &pluralityInfo{
			CurrentNumber: tracker.current + 1,
			TotalCount:    totalBios,
		}

		var nd nameDriver
		if bios[0].SysMemberID != "" { // account bios will have a blank system member ID
			nd = newSystemName(bios[0].SysMemberID, plurality)
		} else {
			nd = newAccountName(targetUserId, ctx.Message.GuildID, plurality, ctx.Session)
		}

		e, err := b.formBioEmbed(nd, bios[0].BioData)
		if err != nil {
			return err
		}

		sentMessage, err := ctx.SendMessageEmbed(ctx.Message.ChannelID, e)
		if err != nil {
			return err
		}

		b.trackerLock.Lock()
		b.trackedEmbeds[sentMessage.ID] = tracker
		b.trackerLock.Unlock()

		for _, v := range []string{previousBioReaction, nextBioReaction} {
			err := ctx.Session.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, v)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

