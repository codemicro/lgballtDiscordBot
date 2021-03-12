package tools

import (
	"github.com/bwmarrin/discordgo"
)

func GetNickname(session *discordgo.Session, userId, guildId string) (string, *discordgo.User, error) {

	member, err := session.GuildMember(guildId, userId)

	if err != nil {

		switch e := err.(type) {
		case *discordgo.RESTError:
			if e.Response.StatusCode == 404 {
				// Can't get the member, so just get the user instead
				user, err := session.User(userId)
				if err != nil {
					return "", nil, err
				}
				return user.Username, user, nil
			}
		default:
			return "", nil, err
		}

	} else {

		name := member.Nick
		if name == "" {
			name = member.User.Username
		}
		return name, member.User, nil

	}
	return "", nil, nil
}
