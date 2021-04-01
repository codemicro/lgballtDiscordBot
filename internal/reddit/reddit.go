package reddit

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

func Start(state *state.State) error {

	ts, err := discordgo.New()
	if err != nil {
		return err
	}

	for _, sub := range config.RedditFeeds {
		go subMonitorSequencer(state, sub, ts)
	}

	return nil
}
