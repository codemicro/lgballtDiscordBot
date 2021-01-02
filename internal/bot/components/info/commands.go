package info

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/go-ping/ping"
	"github.com/hashicorp/go-multierror"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"runtime"
	"time"
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

var rainbowColours = [...]int{
	16711680,
	16753920,
	16776960,
	32768,
	255,
	4915330,
	15631086,
	0,
}

func (i *Info) Info(_ []string, m *harmony.Message) error {
	sinceStart := time.Since(buildInfo.StartTime)

	emb := embed.New().
		Title(fmt.Sprintf("LGBallT bot v%s", buildInfo.Version)).
		Fields(
			embed.NewField().Name("Build date and time").Value(buildInfo.BuildDate).Build(),
			embed.NewField().Name("Go version").Value(buildInfo.GoVersion).Build(),
			embed.NewField().Name("Lines of code").Value(fmt.Sprintf("The bot is currently being powered by %s lines of code, spread across %s files.", buildInfo.LinesOfCode, buildInfo.NumFiles)).Build(),
			embed.NewField().Name("Uptime").Value(fmt.Sprintf("%.0f hours, %.0f minutes and %.0f seconds since start", sinceStart.Hours(), sinceStart.Minutes(), sinceStart.Seconds())).Build(),
	).Build()

	emb.Color = rainbowColours[0]

	msg, err := i.b.SendEmbed(m.ChannelID, emb)
	if err != nil {
		return err
	}

	go func() {
		for c := 1; c < len(rainbowColours); c += 1 {
			time.Sleep(time.Millisecond * 800)
			emb.Color = rainbowColours[c]
			_, err := i.b.Client.Channel(msg.ChannelID).EditEmbed(context.Background(), msg.ID, "", emb)
			if err != nil {
				logging.Error(err, "info colours unable to update")
				return
			}
		}
	}()

	return nil
}