package misc

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
)

func (s *Misc) Avatar(ctx *route.MessageContext) error {
	// Syntax: <user ID>

	id := ctx.Arguments["userId"].(string)

	// get user
	user, err := ctx.Session.User(id)
	if err != nil {
		switch n := err.(type) {
		case *discordgo.RESTError:
			if n.Response.StatusCode == 404 || n.Response.StatusCode == 400 {
				_, err := ctx.SendMessageString(ctx.Message.ChannelID, "⚠ This user doesn't exist.")
				return err
			}
		}
		return err
	}

	// send message
	_, err = ctx.SendMessageString(ctx.Message.ChannelID, user.AvatarURL(""))
	return err
}
