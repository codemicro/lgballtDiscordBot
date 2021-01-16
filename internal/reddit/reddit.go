package reddit

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
)

func Start(state *state.State) error {

	for _, sub := range config.RedditFeeds {
		go subMonitorSequencer(state, sub)
	}

	return nil
}
