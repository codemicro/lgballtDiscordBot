package bios

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"strings"
)

func (b *Bios) RouteMessage(args []string, m *harmony.Message) {

	argLen := len(args)

	if argLen == 0 {
		// This is someone requesting their own bio
		err := b.ReadBio([]string{m.Author.ID}, m)
		if err != nil {
			logging.Error(err)
		}
	} else if argLen == 1 {
		// This can be any one of
		//  * general help text
		//  * system help text
		//  * a user requesting another user bio
		//  * own bio field removal

		if strings.EqualFold(args[0], "help") {
			err := b.Help(args, m)
			if err != nil {
				logging.Error(err, "bios.Bios.Help")
			}

		} else if strings.EqualFold(args[0], "syshelp") {
			err := b.HelpSystem(args, m)
			if err != nil {
				logging.Error(err, "bios.Bios.HelpSystem")
			}

		} else if _, isFieldName := b.validateFieldName(args[0]); isFieldName {
			// This is someone trying to clear a bio field
			err := b.ClearField(args, m)
			if err != nil {
				logging.Error(err, "bios.Bios.ClearField")
			}

		} else {
			// This is someone trying to get the bio of another user
			err := b.ReadBio(args, m)
			if err != nil {
				logging.Error(err, "bios.Bios.ReadBio")
			}
		}
	} else if argLen >= 2 {
		// This can be any one of
		//  * normal field update
		//  * sysmate field update
		//  * sysmate field clear

		accBios, err := db.GetBiosForAccount(m.Author.ID)
		if err != nil {
			logging.Error(err, "db.GetBiosForAccount in bios.RouteMessage")
			return
		}

		if _, isFieldName := b.validateFieldName(args[0]); isFieldName {
			err := b.SetField(args, m)
			if err != nil {
				logging.Error(err, "bios.Bios.SetField")
			}
		} else if tools.IsStringInSlice(args[0], filterForMemberIds(accBios)) { // if member id found for this user

			// Any functions called here *must* validate field names

			if argLen == 2 {
				err := b.ClearFieldSystem(args, m)
				if err != nil {
					logging.Error(err, "bios.Bios.ClearFieldSystem")
				}
			} else {
				// This is assuming that because the first arg is not a valid field name, it is therefore a PK member ID
				err := b.SetFieldSystem(args, m)
				if err != nil {
					logging.Error(err, "bios.Bios.SetFieldSystem")
				}
			}

		} else {
			// This either means that the field name was not recognised or the member ID was not found
			_, err := b.b.SendMessage(m.ChannelID, "That's not a valid field name or a recognised PluralKit member " +
				"ID. Available field names are " + strings.Join(b.data.Fields, ", "))
			if err != nil {
				logging.Error(err, "final parse fail message in bios.RouteMessage")
			}
		}
	}

}

func (b *Bios) RouteAdminMessage(args []string, m *harmony.Message) {

	// TODO: This is going to be royally broken by bios for systems

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

// filterForMemberIds takes a slice of user bios and returns a string slice of all member IDs present in the input
func filterForMemberIds(bsi []db.UserBio) []string {
	var o []string
	for _, bx := range bsi {
		if bx.SysMemberID != "" {
			o = append(o, bx.SysMemberID)
		}
	}
	return o
}
