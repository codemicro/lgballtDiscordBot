package chatchart

import (
	"bytes"
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"github.com/wcharczuk/go-chart/v2"
	"sort"
)

const (
	maxMessages = 5000
	perGroup    = 100
	percentThreshold = 1
)

type percentageWithLabel struct {
	Label      string
	Percentage float64
}

func (c *ChatChart) collectMessages(intent collectionIntent) {

	channelId := intent.ChannelId
	m := intent.Message

	// get channel
	channel := c.b.Client.Channel(channelId)

	// get last message ID
	msg, err := c.b.SendMessage(channelId, "Chatchart collection in progress...")
	if err != nil {
		logging.Error(err, "could not send message from collectMessages")
		return
	}

	lastMessage := msg.ID
	messageCount := 0

	messageUserCount := make(map[string]int)

	// get messages and count per user
	for messageCount < maxMessages {

		messages, err := channel.Messages(context.Background(), "<"+lastMessage, perGroup)

		if err != nil {
			logging.Error(err, "unable to fetch messages for chatchart command")
			return
		}

		if len(messages) == 0 {
			break
		}

		lmt := messages[0].Timestamp

		for _, message := range messages {
			authorName := fmt.Sprintf("%s#%s", message.Author.Username, message.Author.Discriminator)
			v := messageUserCount[authorName]
			v += 1
			messageUserCount[authorName] = v

			if lmt.After(message.Timestamp) {
				lmt = message.Timestamp
				lastMessage = message.ID
			}
		}

		messageCount += len(messages)
	}

	// delete own message
	err = channel.DeleteMessage(context.Background(), msg.ID)
	if err != nil {
		logging.Error(err, "unable to delete own ID collection messages from collectMessages")
		return
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
	// messageUserPercentages := make(map[string]float64)
	for user, val := range messageUserCount {
		percentage := (float64(val) / float64(messageCount)) * 100
		messageUserPercentages = append(messageUserPercentages, percentageWithLabel{
			Label:      user,
			Percentage: percentage,
		})
	}

	sort.Slice(messageUserPercentages, func(i, j int) bool {
		return messageUserPercentages[i].Percentage > messageUserPercentages[j].Percentage
	})

	// make graph

	var chartValues []chart.Value
	for key, val := range messageUserCount {
		chartValues = append(chartValues, chart.Value{
			Label: key,
			Value: float64(val),
		})
	}

	// get channel
	crx, err := channel.Get(context.Background())
	if err != nil {
		logging.Error(err, "failed to get source channel info")
	}

	pie := chart.PieChart{
		Width:  1024,
		Height: 1024,
		Values: chartValues,
	}

	buffer := &tools.ClosingBuffer{Buffer: bytes.NewBuffer([]byte{})}
	err = pie.Render(chart.PNG, buffer)
	if err != nil {
		logging.Error(err, "failed to render chart in collectMessages")
		return
	}

	// form an embed
	emb := embed.Embed{
		Type:        "rich",
		Image:       embed.NewImage("attachment://chart.png"),
		Description: fmt.Sprintf("Chatchart for %s\n", tools.MakeChannelMention(crx.ID)),
	}

	// add fields to embed
	for _, des := range messageUserPercentages {
		emb.Description += fmt.Sprintf("\n**%s**: %.2f%%", des.Label, des.Percentage)
	}

	// send image to user with ping

	_, err = c.b.Client.Channel(m.ChannelID).Send(context.Background(),
		harmony.WithContent(tools.MakePing(m.Author.ID)),
		harmony.WithEmbed(&emb),
		harmony.WithFiles(harmony.FileFromReadCloser(buffer, "chart.png")))

	if err != nil {
		logging.Error(err, "unable to send final chatchart message from collectMessages")
		return
	}
}
