package bios

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
	"strconv"
	"sync"
	"time"
)

func (b *Bios) ReadBio(ctx *route.MessageContext) error {

	targetUserId := ctx.Arguments["user"].(string)

	bios, err := db.GetBiosForAccount(targetUserId)
	if err != nil {
		return err
	}

	if len(bios) == 0 {
		err = ctx.SendErrorMessage("This user hasn't created a bio, or just plain doesn't exist.")
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

		// Present the initial selection box to allow the user to select a bio from multiple choices

		// Find the system ID from one of the bios
		var systemID string
		for _, x := range bios {
			if x.SystemID != "" {
				systemID = x.SystemID
				break
			}
		}

		var warning string
		// Fetch detailed system member info
		members, err := pluralkit.MembersBySystemId(systemID)
		if err != nil {
			if errors.Is(err, pluralkit.ErrorMemberListPrivate) {
				warning = "\n⚠ Cannot retrieve system member names - member list is private"
			} else {
				return err
			}
		}

		acc, err := ctx.Session.User(targetUserId)
		if err != nil {
			return err
		}

		// make strings containing all the account names
		var names []string
		for _, bio := range bios {
			if bio.SysMemberID != "" {
				if systemMember := members.Get(bio.SysMemberID); systemMember != nil {
					name := systemMember.Name
					if systemMember.Nickname != "" {
						name = systemMember.Nickname
					}

					names = append(names, fmt.Sprintf("%s (`%s`)", name, systemMember.Id))
				} else {
					names = append(names, fmt.Sprintf("`%s`", bio.SysMemberID))
				}
			} else {
				name, _, err := common.GetNickname(ctx.Session, acc.ID, ctx.Message.GuildID)
				if err != nil {
					return err
				}
				names = append(names, fmt.Sprintf("Account bio (%s)", name))
			}
		}

		// form one string of all the account names and send as an embed
		var bioSelectionText string
		for i, name := range names {
			bioSelectionText += fmt.Sprintf("**%d** - %s\n", i+1, name)
		}
		m, err := ctx.SendMessageEmbed(ctx.Message.ChannelID, &discordgo.MessageEmbed{
			Description: bioSelectionText + warning,
			Footer: &discordgo.MessageEmbedFooter{
				Text: `Send another message with the number of the bio you'd like to view - for example, "2"` +
					"\nYou will still be able to view other bios afterwards.",
			},
		})
		if err != nil {
			return err
		}

		handlerOnce := new(sync.Once)
		var handlerNumber int
		handlerNumber = ctx.Kit.AddTemporaryMessageHandler(func(ctxb *route.MessageContext) error {
			if ctxb.Message.Author.ID != ctx.Message.Author.ID { // if it's a different person responding
				return nil
			}
			if ctxb.Message.ChannelID != ctx.Message.ChannelID { // if it's in a different channel
				return nil
			}

			// This whole thing is wrapped using sync.Once in order to prevent it being run multiple times if a user
			// sends two messages in very quick succession. That could potentially lead to the message handler a given
			// ID being destroyed multiple times, and removing other unrelated temporary handlers.
			var err error
			handlerOnce.Do(func() {
				defer func() {
					_ = ctx.Session.ChannelMessageDelete(ctxb.Message.ChannelID, ctxb.Message.ID)
					go ctx.Kit.RemoveTemporaryMessageHandler(handlerNumber)
				}()

				selectedNumber, err := strconv.Atoi(ctxb.Message.Content)
				if err != nil {
					err = ctx.SendErrorMessage(fmt.Sprintf("invalid number (%s)", err.Error()))
					if err != nil {
						return
					}
					err = ctx.Session.ChannelMessageDelete(m.ChannelID, m.ID)
					return
				}

				if selectedNumber < 1 || selectedNumber > len(bios) {
					err = ctx.SendErrorMessage(fmt.Sprintf("selected number out of range (min 1, max %d)", len(bios)))
					if err != nil {
						return
					}
					err = ctx.Session.ChannelMessageDelete(m.ChannelID, m.ID)
					return
				}

				tracker := &trackedEmbed{
					accountId:      targetUserId,
					channelId:      ctx.Message.ChannelID,
					bios:           bios,
					timeoutAt:      time.Now().Add(bioTimeoutDuration),
					requestingUser: ctx.Message.Author.ID,
				}

				// send first bio

				plurality := &pluralityInfo{
					CurrentNumber: selectedNumber, // CurrentNumber is only used to show to the user, hence is +1 compared to the target index
					TotalCount:    totalBios,
				}

				selectedNumber -= 1

				var nd nameDriver
				if bios[selectedNumber].SysMemberID != "" { // account bios will have a blank system member ID
					nd = newSystemName(bios[selectedNumber].SysMemberID, plurality)
				} else {
					nd = newAccountName(targetUserId, ctx.Message.GuildID, plurality, ctx.Session)
				}

				var e *discordgo.MessageEmbed
				e, err = b.formBioEmbed(nd, bios[selectedNumber].BioData)
				if err != nil {
					return
				}

				var sentMessage *discordgo.Message
				sentMessage, err = ctx.Session.ChannelMessageEditEmbed(ctx.Message.ChannelID, m.ID, e)
				if err != nil {
					return
				}

				b.trackerLock.Lock()
				b.trackedEmbeds[sentMessage.ID] = tracker
				b.trackerLock.Unlock()

				for _, v := range []string{previousBioReaction, nextBioReaction} {
					err = ctx.Session.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, v)
					if err != nil {
						return
					}
				}

				return
			})
			return err
		})
	}

	return nil
}
