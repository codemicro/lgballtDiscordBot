package reddit

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"math/rand"
	"time"
)

func Start(state *state.State) error {

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
