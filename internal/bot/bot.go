package bot

import (
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/bios"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/info"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/roles"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"strings"
)

const partyRoleId = "698570587567685703"

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

		// TODO: Command parsing, but actually make it decent this time

		// Remove prefix and split by spaces
		messageComponents := strings.Split(
			strings.TrimPrefix(m.Content, b.Prefix),
			" ")

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

			// ---------- BIOS ----------

			bioComponents := messageComponents[1:]

			if len(bioComponents) == 0 {
				// This is someone requesting their own bio
				err := bioComponent.ReadBio([]string{m.Author.ID}, m)
				if err != nil {
					logging.Error(err)
				}
			} else if len(bioComponents) == 1 {

				// This is either someone clearing a field from their bio OR someone requesting a user's bio by ID
				// OR someone requesting the help menu

				if strings.EqualFold(bioComponents[0], "help") {
					err := bioComponent.Help(bioComponents, m)
					if err != nil {
						logging.Error(err)
					}
				} else if _, isFieldName := bioComponent.ValidateFieldName(bioComponents[0]); isFieldName {
					// This is someone trying to clear a bio field
					err := bioComponent.ClearField(bioComponents, m)
					if err != nil {
						logging.Error(err)
					}
				} else {
					// This is someone trying to get the bio of another user
					err := bioComponent.ReadBio(bioComponents, m)
					if err != nil {
						logging.Error(err)
					}
				}
			} else {
				// This is triggered where there are two or more arguments
				// The only command that satisfies this is the set bio field command

				err := bioComponent.SetField(bioComponents, m)
				if err != nil {
					logging.Error(err)
				}
			}

		} else if strings.EqualFold(messageComponents[0], "biof") {

			if m.Author.ID == "289130374204751873" {  // 0x5444#8669
				adminBioComponents := messageComponents[1:]

				if len(adminBioComponents) == 1 {
					// This is me getting the value of a user's bio
					err := bioComponent.AdminReadRawBio(adminBioComponents, m)
					if err != nil {
						logging.Error(err)
					}
				} else {
					// This is only triggered for two or more arguments
					err := bioComponent.AdminSetRawBio(adminBioComponents, m)
					if err != nil {
						logging.Error(err)
					}
				}

			}

		} else if strings.EqualFold(messageComponents[0], "info") {

			// ---------- INFO ----------

			infoComponents := messageComponents[1:]

			if len(infoComponents) == 1 {
				if strings.EqualFold(infoComponents[0], "ping") {
					err := infoComponent.Ping(infoComponents, m)
					if err != nil {
						logging.Error(err)
					}
				}
			}

		} else if strings.EqualFold(messageComponents[0], "roles") {

			// ---------- ROLES ----------

			if tools.IsStringInSlice(partyRoleId, m.Member.Roles) {

				if len(messageComponents) >= 2 {  // required minimum of two arguments
					instruction := messageComponents[1]
					roleComponents := messageComponents[2:]

					if strings.EqualFold(instruction, "track") && len(roleComponents) >= 3 {
						err := roleComponent.TrackReaction(roleComponents, m)
						if err != nil {
							logging.Error(err)
						}
					} else if strings.EqualFold(instruction, "untrack") && len(roleComponents) >= 2 {
						err := roleComponent.UntrackReaction(roleComponents, m)
						if err != nil {
							logging.Error(err)
						}
					}

				}

			}
			
		} else if strings.EqualFold(messageComponents[0], "broken") {
			err := infoComponent.Broken([]string{}, m)
			if err != nil {
				logging.Error(err)
			}
		}

	})

	b.Client.OnMessageReactionAdd(func (r *harmony.MessageReaction) {
		err := roleComponent.ReactionAdd(r)
		if err != nil {
			logging.Error(err)
		}
	})

	b.Client.OnMessageReactionRemove(func (r *harmony.MessageReaction) {
		err := roleComponent.ReactionRemove(r)
		if err != nil {
			logging.Error(err)
		}
	})

	return nil
}
