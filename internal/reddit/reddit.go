package reddit

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"math/rand"
	"time"
)

func Start(state *state.State) error {

	if config.DebugMode {
		// Don't run Reddit feeds if we're working in debug mode to prevent accidentally sending rather a lot of
		// requests to Reddit in a short space of time. They tend not to like that.
		return nil
	}

	ts, err := discordgo.New()
	if err != nil {
		return err
	}

	for _, sub := range config.RedditFeeds {
		go func() {
			// Stagger timer starts - Reddit only likes you to send at most 1 request every two seconds.
			// This is hacky way to try and ensure that happens
			time.Sleep(time.Second * time.Duration(rand.Intn(28) + 2))
			go subMonitorSequencer(state, sub, ts)
		}()
	}

	return nil
}
