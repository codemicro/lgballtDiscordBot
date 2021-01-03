package reddit

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/codemicro/dishook"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/mmcdole/gofeed"
	"strings"
	"time"
)

func subMonitorSequencer(state *tools.State, info config.RedditFeedInfo) {

	var idCache []string

	state.AddGoroutine()

	ticker := time.NewTicker(time.Duration(info.Interval) * time.Minute)
	// ticker := time.NewTicker(time.Second * 15)
	finished := make(chan bool)

	go func() {
		state.WaitUntilShutdownTrigger()
		ticker.Stop()
		finished <- true
	}()

	var jumpOut bool
	for {

		if jumpOut {
			break
		}

		// So... apparently it's possible to break out of a select statement
		// Half an hour wasted debugging shutdown deadlocks... :)
		select {
		case <-finished:
			jumpOut = true
		case <-ticker.C:
			subMonitorAction(info, &idCache)
		}
	}

	state.FinishGoroutine()

}

func subMonitorAction(info config.RedditFeedInfo, idCache *[]string) {

	fmt.Println("action run")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(info.RssUrl, ctx)
	if err != nil {
		logging.Error(err, fmt.Sprintf("failed to fetch RSS feed from URL %s within specified timeout",
			info.RssUrl))
		return
	}

	// This is used for when the action runs the first time so old messages aren't posted like forty times whenever the
	// bot starts
	sendMessages := true
	if len(*idCache) == 0 {
		sendMessages = false
	}

	var newItems []*gofeed.Item

	for _, item := range feed.Items {
		if !tools.IsStringInSlice(item.GUID, *idCache) {
			newItems = append(newItems, item)
			*idCache = append(*idCache, item.GUID)
		}
	}

	if sendMessages {

		// for each new item, form a webhook message and fire that off with incredible speed and vigour
		// what am I talking about
		// just send it

		for _, redditPost := range newItems {
			wh := dishook.NewMessage()

			var emb dishook.Embed
			emb.Title = redditPost.Title
			emb.Author = dishook.EmbedAuthor{
				Name: redditPost.Author.Name,
				URL:  fmt.Sprintf("https://www.reddit.com%s", redditPost.Author.Name),
			}
			emb.URL = redditPost.Link
			emb.Footer = dishook.EmbedFooter{
				Text: fmt.Sprintf("New post at %s - /r/%s",
					redditPost.UpdatedParsed.Format("Mon Jan 2 15:04:05 MST 2006"), redditPost.Categories[0]),
			}

			extVals := redditPost.Extensions["media"]["thumbnail"]
			if len(extVals) >= 1 {
				v, ok := extVals[0].Attrs["url"]
				if ok {
					emb.Thumbnail = dishook.EmbedImage{URL: v}
				}
			}

			content := redditPost.Content
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
			if err != nil {
				logging.Error(err, "goquery initialisation failed")
				return
			}
			content = doc.Text()

			if len(content) > 70 {
				content = content[:70] + "..."
			}

			emb.Description = content

			wh.Embeds = append(wh.Embeds, emb)
			wh.Username = "/r/" + redditPost.Categories[0]
			wh.AvatarURL = info.IconUrl
			_, err = wh.Send(info.Webhook, true)
			if err != nil {
				logging.Error(err, fmt.Sprintf("unable to send to webhook URL %s", info.Webhook))
				return
			}

		}
	}

}