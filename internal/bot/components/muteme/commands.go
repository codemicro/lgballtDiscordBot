package muteme

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"time"
)

func (mm *MuteMe) Trigger(ctx *route.MessageContext) error {

	duration := ctx.Arguments["duration"].(time.Duration)
	finishTime := time.Now().Add(duration)
	finishTimeString := finishTime.UTC().Format(time.RFC822)

	if finishTime.After(time.Now().Add(time.Hour * 24 * 366)) {
		_, err := ctx.SendMessageString(ctx.Message.ChannelID, "Maximum duration is 365 days")
		return err
	}

	emb := &discordgo.MessageEmbed{
		Title: "Confirmation",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "React to this message with ✅ if you wish to mute yourself. (❌ to cancel)",
		},
		Description: fmt.Sprintf(muteMeText, finishTimeString),
	}

	err := ctx.Kit.NewConfirmation(
		ctx.Message.ChannelID,
		ctx.Message.Author.ID,
		emb,
		func(ctx *route.ReactionContext) error {
			// do mute

			// get user
			user, err := ctx.Session.GuildMember(ctx.Reaction.GuildID, ctx.Reaction.UserID)
			if err != nil {
				return err
			}

			// determine which roles need to be removed
			var rolesToRemove []string
			for _, role := range config.MuteMe.RolesToRemove {
				if tools.IsStringInSlice(role, user.Roles) {
					rolesToRemove = append(rolesToRemove, role)
				}
			}

			// add DB entry
			um := db.UserMute{
				UserId:       ctx.Reaction.UserID,
				GuildId:      ctx.Reaction.GuildID,
				EndTime:      finishTime.Unix(),
				RemovedRoles: rolesToRemove,
			}
			err = um.Create()
			if err != nil {
				return err
			}

			// add timeout role
			err = ctx.Session.GuildMemberRoleAdd(ctx.Reaction.GuildID, ctx.Reaction.UserID, config.MuteMe.TimeoutRole)
			if err != nil {
				return err
			}

			// remove roles
			for _, role := range rolesToRemove {
				err = ctx.Session.GuildMemberRoleRemove(ctx.Reaction.GuildID, ctx.Reaction.UserID, role)
				if err != nil {
					return err
				}
			}

			_, _ = ctx.Session.ChannelMessageEdit(ctx.Reaction.ChannelID, ctx.Reaction.MessageID, "Muted!")
			return nil
		},
		func(ctx *route.ReactionContext) error {
			_, _ = ctx.Session.ChannelMessageEdit(ctx.Reaction.ChannelID, ctx.Reaction.MessageID, "Cancelled!")
			return nil
		},
	)

	return err

}
