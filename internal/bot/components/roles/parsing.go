package roles

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"strings"
)

func (r *Roles) RouteMessage(args []string, m *harmony.Message) {

	if tools.IsStringInSlice(config.AdminRole, m.Member.Roles) || config.DebugMode {

		if len(args) >= 1 { // required minimum of two arguments
			instruction := args[0]
			roleComponents := args[1:]

			if strings.EqualFold(instruction, "track") && len(roleComponents) >= 3 {
				err := r.TrackReaction(roleComponents, m)
				if err != nil {
					logging.Error(err)
				}
			} else if strings.EqualFold(instruction, "untrack") && len(roleComponents) >= 2 {
				err := r.UntrackReaction(roleComponents, m)
				if err != nil {
					logging.Error(err)
				}
			}

		}

	}

}
