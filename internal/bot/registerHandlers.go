package bot

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/bios"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/chatchart"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/info"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/misc"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/roles"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/verification"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	harmonyChannel "github.com/skwair/harmony/channel"
	"strings"
)

func RegisterHandlers(b *core.Bot) error {
	// Load components
	bioComponent, err := bios.New(b)
	if err != nil {
		return err
	}

	infoComponent, err := info.New(b)
	if err != nil {
		return err
	}

	roleComponent, err := roles.New(b)
	if err != nil {
		return err
	}

	chartComponent, err := chatchart.New(b)
	if err != nil {
		return err
	}

	miscComponent, err := misc.New(b)
	if err != nil {
		return err
	}

	verificationComponent, err := verification.New(b)
	if err != nil {
		return err
	}

	// Register handler functions
	b.Client.OnMessageCreate(func(m *harmony.Message) {

		// Ignore own messages
		if isSelf, err := b.IsSelf(m.Author.ID); err != nil {
			logging.Error(err)
			return
		} else if isSelf {
			return
		}

		// ignore bots
		if m.Author.Bot {
			return
		}

		if !strings.HasPrefix(m.Content, b.Prefix) {
			return
		}

		// ignore DMs and group DMs
		channel, err := b.Client.Channel(m.ChannelID).Get(context.Background())
		if err != nil {
			logging.Error(err)
			return
		}
		if channel.Type == harmonyChannel.TypeDM || channel.Type == harmonyChannel.TypeGroupDM {
			return
		}

		// Remove prefix and split by spaces
		messageComponents := strings.Split(
			strings.TrimPrefix(m.Content, b.Prefix),
			" ")

		// Replace at most one newline with a space so <p>cmd\nwords is registered as <p>cmd words
		// This is based off the first item in the message when split by spaces only to prevent newlines being randomly
		// replaced in the middle of a bio field, for example
		modMessageComp := strings.Split(strings.Replace(messageComponents[0], "\n", " ", 1), " ")
		if len(modMessageComp) > 1 {
			messageComponents = append(modMessageComp, messageComponents[1:]...)
		}


		// strings.Split will never return a empty slice - this can lead to a slice with a single empty string in it
		// being returned, signifying that the input string was the prefix alone.
		// in this case, we just empty the message components slice
		if len(messageComponents) == 1 {
			if messageComponents[0] == "" {
				messageComponents = []string{}
			}
		}

		// There's nothing in the command, so we should ignore it
		if len(messageComponents) < 1 {
			return
		}

		if strings.EqualFold(messageComponents[0], "bio") {
			bioComponent.RouteMessage(messageComponents[1:], m)

		} else if strings.EqualFold(messageComponents[0], "biof") {
			bioComponent.RouteAdminMessage(messageComponents[1:], m)

		} else if strings.EqualFold(messageComponents[0], "info") {
			infoComponent.RouteMessage(messageComponents[1:], m)

		} else if strings.EqualFold(messageComponents[0], "roles") {
			roleComponent.RouteMessage(messageComponents[1:], m)

		} else if strings.EqualFold(messageComponents[0], "chatchart") {
			chartComponent.RouteMessage(messageComponents[1:], m)

		} else if strings.EqualFold(messageComponents[0], "avatar") {

			err := miscComponent.Avatar(messageComponents[1:], m)
			if err != nil {
				logging.Error(err, "miscComponent.Avatar")
			}

		} else if strings.EqualFold(messageComponents[0], "emoji") {

			if len(messageComponents[1:]) >= 1 {
				err := miscComponent.Emoji(messageComponents[1:], m)
				if err != nil {
					logging.Error(err, "miscComponent.Emoji")
				}
			}

		} else if strings.EqualFold(messageComponents[0], "verify") && (m.ChannelID == verification.InputChannelId ||
			config.DebugMode) {

			// ---------- VERIFICATION -------------

			err := verificationComponent.Verify(messageComponents[1:], m, true)
			if err != nil {
				logging.Error(err, "verificationComponent.Verify")
			}

		} else if strings.EqualFold(messageComponents[0], "verifyf") &&
			(m.ChannelID == verification.InputChannelId || config.DebugMode) &&
			(tools.IsStringInSlice(config.AdminRole, m.Member.Roles) || config.DebugMode) {

			// ---------- FORCE VERIFY -------------

			err := verificationComponent.FVerify(messageComponents[1:], m)
			if err != nil {
				logging.Error(err, "verificationComponent.FVerify")
			}

		} else if (strings.EqualFold("ban", messageComponents[0]) ||
			strings.EqualFold("kick", messageComponents[0])) &&
			(tools.IsStringInSlice(config.AdminRole, m.Member.Roles) || config.DebugMode) {

			// ---------- KICK/BAN TRIGGER -------------

			err := verificationComponent.RecordRemoval(messageComponents, m)
			if err != nil {
				logging.Error(err, "verificationComponent.RecordRemoval")
			}

		} else if strings.EqualFold("pressf", messageComponents[0]) {
			// ---------- PRESSF -------------

			err := miscComponent.PressF(messageComponents[1:], m)
			if err != nil {
				logging.Error(err, "miscComponent.PressF")
			}

		} else if strings.EqualFold("shutdown", messageComponents[0]) &&
			(tools.IsStringInSlice(config.AdminRole, m.Member.Roles) || m.Author.ID == config.OwnerId ||
				config.DebugMode) {

			// Goodnight.

			logging.Info(fmt.Sprintf("Shutting down from command by %s %s#%s", m.Author.ID, m.Author.Username,
				m.Author.Discriminator))

			_, _ = b.SendMessage(m.ChannelID, "Bye-bye!")

			b.State.ShutdownSignal <- state.CustomShutdownSignal{}

		}

	})

	b.Client.OnMessageReactionAdd(func(r *harmony.MessageReaction) {

		// Ignore own messages
		if isSelf, err := b.IsSelf(r.UserID); err != nil {
			logging.Error(err)
			return
		} else if isSelf {
			return
		}

		if r.ChannelID == verification.OutputChannelId {
			err := verificationComponent.AdminDecision(r)
			if err != nil {
				logging.Error(err, "verificationComponent.AdminDecision")
			}
		}

		err := bioComponent.ReactionAdd(r)
		if err != nil {
			logging.Error(err, "bioComponent.ReactionAdd")
		}

		err = roleComponent.ReactionAdd(r)
		if err != nil {
			logging.Error(err)
		}

		err = miscComponent.PressFReaction(r)
		if err != nil {
			logging.Error(err, "miscComponent.PressFReaction")
		}
	})

	b.Client.OnMessageReactionRemove(func(r *harmony.MessageReaction) {

		// Ignore own messages
		if isSelf, err := b.IsSelf(r.UserID); err != nil {
			logging.Error(err)
			return
		} else if isSelf {
			return
		}

		err := roleComponent.ReactionRemove(r)
		if err != nil {
			logging.Error(err)
		}
	})

	return nil
}
