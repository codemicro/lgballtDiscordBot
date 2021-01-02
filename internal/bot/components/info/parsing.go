package info

import (
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/skwair/harmony"
	"strings"
)

func (i *Info) RouteMessage(args []string, m *harmony.Message) {

	if len(args) == 1 {
		if strings.EqualFold(args[0], "ping") {
			err := i.Ping(args, m)
			if err != nil {
				logging.Error(err)
			}
		}
	} else {
		err := i.Info(args, m)
		if err != nil {
			logging.Error(err, "info.Info.Info")
		}
	}

}
