package chatchart

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/wcharczuk/go-chart/v2"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	maxMessages      = 5000
	percentThreshold = 1
	usernameMaxLen   = 16
)

type percentageWithLabel struct {
	Label      string
	Percentage float64
}

func magicalUsernameTrim(in string) string {
	// magic!!1!
	x := strings.Split(in, "#")
	ix := usernameMaxLen
	var app string
	if len(x[0]) < ix {
		ix = len(x[0])
	} else {
		app = "..."
	}
	o := x[0][:ix] + app
	if len(x) >= 2 {
		return o + "#" + x[1]
	}
	return o
}

const discordEpoch = 1420070400000

func makeSnowflakeAtTime(t time.Time) string {
	x := (uint64(t.UnixNano()/1000000) - discordEpoch) << 22
	return strconv.FormatUint(x, 10)
}

func (c *ChatChart) collectMessages(intent collectionIntent) {

	// TODO: zerolog

	channelId := intent.ChannelId
	ctx := intent.Ctx

	lastMessage := makeSnowflakeAtTime(time.Now().Add(time.Minute * 5))
	messageCount := 0

	messageUserCount := make(map[string]int)

	// get messages and count per user
	for messageCount < maxMessages {

		// trigger typing
		err := ctx.Session.ChannelTyping(channelId)
		if err != nil {
			logging.Warn(err.Error())
		}

		messages, err := ctx.Session.ChannelMessages(channelId, 100, lastMessage, "", "")

		if err != nil {
			logging.Error(err, "unable to fetch messages for chatchart command")
			return
		}

		if len(messages) == 0 {
			break
		}

		lmt, _ := discordgo.SnowflakeTimestamp(messages[0].ID)

		for _, message := range messages {
			authorName := fmt.Sprintf("%s#%s", message.Author.Username, message.Author.Discriminator)
			v := messageUserCount[authorName]
			v += 1
			messageUserCount[authorName] = v

			nt, _ := discordgo.SnowflakeTimestamp(message.ID)

			if lmt.After(nt) {
				lmt = nt
				lastMessage = message.ID
			}
		}

		messageCount += len(messages)
	}

	// filter out groups that are < 5 percent of the total
	var otherTotal int
	threshold := (messageCount / 100) * percentThreshold
	for user, val := range messageUserCount {
		if val < threshold {
			otherTotal += val
			delete(messageUserCount, user)
		}
	}

	if otherTotal != 0 {
		messageUserCount["Other"] = otherTotal
	}

	// find percentages
	var messageUserPercentages []percentageWithLabel
	var chartValues []chart.Value
	for user, val := range messageUserCount {
		percentage := (float64(val) / float64(messageCount)) * 100
		messageUserPercentages = append(messageUserPercentages, percentageWithLabel{
			Label:      user,
			Percentage: percentage,
		})
		chartValues = append(chartValues, chart.Value{
			Label: magicalUsernameTrim(user),
			Value: percentage,
		})
	}

	sort.Slice(messageUserPercentages, func(i, j int) bool {
		return messageUserPercentages[i].Percentage > messageUserPercentages[j].Percentage
	})

	// get channel
	crx, err := ctx.Session.Channel(channelId)
	if err != nil {
		logging.Error(err, "failed to fetch channel")
	}

	// make graph
	barWidth := 25
	barSpacing := 10

	pie := chart.BarChart{
		Width: 25 + (len(chartValues) * (barWidth + barSpacing)),
		// Height: 1024,
		Background: chart.Style{
			Padding: chart.Box{
				Bottom: 100,
			},
		},
		BarWidth:   25,
		BarSpacing: 10,
		Bars:       chartValues,
		XAxis:      chart.Style{TextRotationDegrees: 90},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = pie.Render(chart.PNG, buffer)
	if err != nil {
		logging.Error(err, "failed to render chart in collectMessages")
		return
	}

	// form an embed
	emb := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Chatchart for #%s", crx.Name),
		Image: &discordgo.MessageEmbedImage{URL: "attachment://chart.png"},
	}

	// add fields to embed
	for _, des := range messageUserPercentages {
		emb.Description += fmt.Sprintf("**%s**: %.2f%%\n", des.Label, des.Percentage)
	}

	emb.Description = strings.TrimSpace(emb.Description)

	// send image to user with ping

	msend := &discordgo.MessageSend{
		Embed:   emb,
		Content: tools.MakePing(ctx.Message.Author.ID),
		Files: []*discordgo.File{
			{Name: "chart.png", Reader: buffer},
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Users: []string{ctx.Message.Author.ID},
		},
		Reference: &discordgo.MessageReference{
			MessageID: ctx.Message.ID,
			ChannelID: ctx.Message.ChannelID,
			GuildID:   ctx.Message.GuildID,
		},
	}

	_, err = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, msend)
	// If the requesting message has been deleted, the message reference will fail and an error will be returned.
	if err != nil {
		// retry without the message reference
		msend.Reference = nil
		// because the buffer has been read, the chart will need to be re-rendered
		err = pie.Render(chart.PNG, buffer)
		if err != nil {
			logging.Error(err, "failed to render chart in collectMessages after the first attempt to send a message")
			return
		}

		_, err = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, msend)
		if err != nil {
			logging.Error(err, "unable to send final chatchart message from collectMessages")
			return
		}
	}
}
