package chatchart

import (
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/skwair/harmony"
)

func (c *ChatChart) RouteMessage(args []string, m *harmony.Message) {
	if len(args) >= 1 {
		err := c.TriggerCollection(args, m)
		if err != nil {
			logging.Error(err)
		}
	}

}
