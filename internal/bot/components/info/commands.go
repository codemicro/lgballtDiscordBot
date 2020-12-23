package info

import (
	"context"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/hashicorp/go-multierror"
	"github.com/skwair/harmony"
	"runtime"
)

func (i *Info) Ping(_ []string, m *harmony.Message) error {

	_ = i.b.Client.Channel(m.ChannelID).TriggerTyping(context.Background())

	pinger, err := ping.NewPinger("www.discord.com")
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" { // See https://github.com/go-ping/ping#windows
		pinger.SetPrivileged(true)
	}

	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		_, mErr := i.b.SendMessage(m.ChannelID, "Unable to complete ping.")
		if mErr != nil {
			err = multierror.Append(err, mErr)
		}
		return err
	}
	stats := pinger.Statistics()
	_, err = i.b.SendMessage(m.ChannelID, fmt.Sprintf("Pong! Average ping time was `%dms`",
		stats.AvgRtt.Milliseconds()))

	return nil

}
