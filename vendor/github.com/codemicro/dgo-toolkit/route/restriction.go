package route

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// CommandRestriction is a function that returns true if a command can be run based on the current state (eg. details
// about the user, roles, channel, etc)
type CommandRestriction func(session *discordgo.Session, message *discordgo.MessageCreate) (bool, error)

func isStringInSlice(needle string, haystack []string) (found bool) {
	for _, v := range haystack {
		if strings.EqualFold(needle, v) {
			found = true
			break
		}
	}
	return
}

// RestrictionByRole creates a CommandRestriction that requires the commanding guild member to have a given role ID
func RestrictionByRole(roleId string) CommandRestriction {
	return func(_ *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
		return isStringInSlice(roleId, message.Member.Roles), nil
	}
}

// RestrictionByChannel creates a CommandRestriction that requires the command to have been sent in a channel with a
// given ID
func RestrictionByChannel(channelIds ...string) CommandRestriction {
	return func(_ *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
		return isStringInSlice(message.ChannelID, channelIds), nil
	}
}
