package info

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"time"
)

func (*Info) Ping(ctx *route.MessageContext) error {

	_, err := ctx.SendMessageString(ctx.Message.ChannelID, fmt.Sprintf("Pong! Current heartbeat latency is "+
		"`%dms`", ctx.Session.HeartbeatLatency().Milliseconds()))

	return err
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

func (i *Info) Info(ctx *route.MessageContext) error {
	sinceStart := time.Since(buildInfo.StartTime)

	hours := int64(sinceStart.Hours())
	minutes := int64(sinceStart.Minutes()) % 60
	seconds := int64(sinceStart.Seconds()) % 60

	earthRotations := float64(hours) / 24

	cmds, reactions := ctx.Kit.GetNums()

	emb := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("LGBallT bot v%s", buildInfo.Version),
		Color: rainbowColours[0],
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Build date and time", Value: buildInfo.BuildDate},
			{Name: "Go version and arch", Value: buildInfo.GoVersion},
			{Name: "Lines of code", Value: fmt.Sprintf("The bot is currently being powered by %s lines of code, spread across %s files.", buildInfo.LinesOfCode, buildInfo.NumFiles)},
			{Name: "Registered commands", Value: fmt.Sprintf("There are currently %d commands and %d reaction handlers registered in the bot.", cmds, reactions)},
			{Name: "Changelog", Value: buildInfo.ChangelogURL},
			{Name: "Uptime", Value: fmt.Sprintf("%d hours, %d minutes and %d seconds since start\nThat's %0.3f rotations of the earth", hours, minutes, seconds, earthRotations)},
		},
	}

	msg, err := ctx.SendMessageEmbed(ctx.Message.ChannelID, emb)
	if err != nil {
		return err
	}

	go func() {
		for c := 1; c < len(rainbowColours); c += 1 {
			time.Sleep(time.Millisecond * 800)
			emb.Color = rainbowColours[c]
			_, err = ctx.Session.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, emb)
			if err != nil {
				logging.Error(err, "info colours unable to update")
				return
			}
		}
	}()

	return nil
}
