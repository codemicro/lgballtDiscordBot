package bios

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/skwair/harmony"
	"strings"
)

func (b *Bios) RouteMessage(args []string, m *harmony.Message) {

	if len(args) == 0 {
		// This is someone requesting their own bio
		err := b.ReadBio([]string{m.Author.ID}, m)
		if err != nil {
			logging.Error(err)
		}
	} else if len(args) == 1 {

		// This is either someone clearing a field from their bio OR someone requesting a user's bio by ID
		// OR someone requesting the help menu

		if strings.EqualFold(args[0], "help") {
			err := b.Help(args, m)
			if err != nil {
				logging.Error(err)
			}
		} else if _, isFieldName := b.validateFieldName(args[0]); isFieldName {
			// This is someone trying to clear a bio field
			err := b.ClearField(args, m)
			if err != nil {
				logging.Error(err)
			}
		} else {
			// This is someone trying to get the bio of another user
			err := b.ReadBio(args, m)
			if err != nil {
				logging.Error(err)
			}
		}
	} else {
		// This is triggered where there are two or more arguments
		// The only command that satisfies this is the set bio field command

		err := b.SetField(args, m)
		if err != nil {
			logging.Error(err)
		}
	}

}

func (b *Bios) RouteAdminMessage(args []string, m *harmony.Message) {

	if m.Author.ID == config.OwnerId {
		if len(args) == 1 {
			// This is me getting the value of a user's bio
			err := b.AdminReadRawBio(args, m)
			if err != nil {
				logging.Error(err)
			}
		} else {
			// This is only triggered for two or more arguments
			err := b.AdminSetRawBio(args, m)
			if err != nil {
				logging.Error(err)
			}
		}
	}
}
