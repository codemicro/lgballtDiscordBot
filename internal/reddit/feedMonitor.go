package reddit

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
	"time"
)

func subMonitorSequencer(state *state.State, info config.RedditFeedInfo, ts *discordgo.Session) {

	var idCache []string

	ticker := time.NewTicker(time.Duration(info.Interval) * time.Minute)
	// ticker = time.NewTicker(time.Second * 10)
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
			subMonitorAction(info, &idCache, ts)
		}
	}

	state.FinishGoroutine()

}

var (
	contentFilterRegex = regexp.MustCompile(`(?m) *submitted +by +\/u\/.+ +\[link\] +\[comments\] *`)
	webhookIdTokenRegex = regexp.MustCompile(`(?m)\/webhooks\/(.+)\/(.+)\/?`)
)

func subMonitorAction(info config.RedditFeedInfo, idCache *[]string, ts *discordgo.Session) {

	logger := log.With().Str("area", "redditSubMonitorAction").Str("rssUrl", info.IconUrl).Logger()

	var hasRestarted bool

ActionStart:
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	fp.UserAgent = fmt.Sprintf("discord:lgballtDiscordBot:v%s (by /u/td_brokeit)", buildInfo.Version)
	feed, err := fp.ParseURLWithContext(info.RssUrl, ctx)
	if err != nil {

		if te, ok := err.(gofeed.HTTPError); ok {
			if te.StatusCode == 429 && !hasRestarted {
				// ratelimit has been hit, for the first time this run

				pauseDuration := 1

				if pauseDuration < info.Interval {
					// If we don't get here, it means that the time we want to wait is less than the interval between
					// feed checks. In that case, we abandon ship here to prevent having multiple request workers
					// running at once.
					time.Sleep(time.Minute * time.Duration(pauseDuration))
					hasRestarted = true
					goto ActionStart // >:) sorry not really sorry
				}

			}
		}

		logger.Error().Err(err).Msgf("failed to fetch RSS feed from URL %s within specified timeout",
			info.RssUrl)
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

			var formattedTime time.Time
			if redditPost.UpdatedParsed == nil {
				formattedTime = time.Now()
			} else {
				formattedTime = *redditPost.UpdatedParsed
			}

			whParams := &discordgo.WebhookParams{
				Username:        "r/" + redditPost.Categories[0],
				AvatarURL:       info.IconUrl,
				Embeds:          []*discordgo.MessageEmbed{
					{
						Title:       redditPost.Title,
						URL: redditPost.Link,
						Footer:      &discordgo.MessageEmbedFooter{
							Text:         fmt.Sprintf("New post at %s - /r/%s",
								formattedTime.Format("Mon Jan 2 15:04:05 MST 2006"), redditPost.Categories[0]),
						},
						Author:      &discordgo.MessageEmbedAuthor{
							URL:          fmt.Sprintf("https://www.reddit.com%s", redditPost.Author.Name),
							Name:         redditPost.Author.Name,
						},
					},
				},
			}


			extVals := redditPost.Extensions["media"]["thumbnail"]
			if len(extVals) >= 1 {
				v, ok := extVals[0].Attrs["url"]
				if ok {
					whParams.Embeds[0].Thumbnail = &discordgo.MessageEmbedThumbnail{URL: v}
				}
			}

			content := redditPost.Content
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
			if err != nil {
				logger.Error().Err(err).Msg("goquery initialisation failed")
				return
			}
			content = contentFilterRegex.ReplaceAllString(doc.Text(), "")

			if len(content) > 70 {
				content = content[:70] + "..."
			}

			whParams.Embeds[0].Description = content

			mtx := webhookIdTokenRegex.FindStringSubmatch(info.Webhook)
			if mtx == nil {
				logger.Error().Err(err).Msg(fmt.Sprintf("unable to extract webhook ID and token from URL %s", info.Webhook))
				return
			}

			_, err = ts.WebhookExecute(
				mtx[1],
				mtx[2],
				true,
				whParams,
			)

			if err != nil {
				logger.Error().Err(err).Msg(fmt.Sprintf("unable to send to webhook URL %s", info.Webhook))
				return
			}

		}
	}

}
