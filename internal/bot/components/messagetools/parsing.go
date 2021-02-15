package messagetools

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"strings"
)

func (t *MessageTools) RouteMessage(args []string, m *harmony.Message) {

	if !(tools.IsStringInSlice(config.AdminRole, m.Member.Roles) || m.Author.ID != config.OwnerId) || len(args) < 2 {
		return
	}

	op, arguments := args[0], args[1:]

	if strings.EqualFold(op, "send") {
		err := t.Send(arguments, m)
		if err != nil {
			logging.Error(err, "messagetools.MessageTools.Send")
		}
	} else if strings.EqualFold(op, "edit") {
		err := t.Edit(arguments, m)
		if err != nil {
			logging.Error(err, "messagetools.MessageTools.Edit")
		}
	}

}
