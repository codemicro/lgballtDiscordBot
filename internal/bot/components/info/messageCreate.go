package info

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/go-ping/ping"
	"github.com/skwair/harmony"
	"strings"
)

func (b *info) onMessageCreate(m *harmony.Message) {

	if isSelf, err := b.b.IsSelf(m.Author.ID); err != nil {
		logging.Error(err)
		return
	} else if isSelf {
		return
	}

	// ignore bots
	if m.Author.Bot {
		return
	}

	fmt.Println("processed here")

	if strings.HasPrefix(m.Content, b.commandPrefix) {
		command := tools.GetCommand(m.Content, b.commandPrefix)

		if len(command) == 1 {
			if strings.EqualFold(command[0], "ping") {

				// Run ping command

				pinger, err := ping.NewPinger("www.discord.com")
				pinger.SetPrivileged(true)
				if err != nil {
					logging.Error(err)
					return
				}
				pinger.Count = 3
				err = pinger.Run() // blocks until finished
				if err != nil {
					logging.Error(err)
					_, err = b.b.SendMessage(m.ChannelID, "Unable to complete ping.")
					if err != nil {
						logging.Error(err)
					}
					return
				}
				stats := pinger.Statistics() // get send/receive/rtt stats
				_, err = b.b.SendMessage(m.ChannelID, fmt.Sprintf("Pong! Average response time was `%dms`", stats.AvgRtt.Milliseconds()))

			}
		}

	}

}