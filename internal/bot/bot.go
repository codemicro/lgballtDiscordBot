package bot

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"time"
)

func Start(state *tools.State) error {

	client, err := harmony.NewClient(config.Token, harmony.WithGatewayIntents(harmony.GatewayIntentUnprivileged))
	if err != nil {
		return err
	}

	b := core.New(client, config.Prefix)
	err = RegisterHandlers(b)
	if err != nil {
		return err
	}

	if err = client.Connect(context.Background()); err != nil {
		return err
	}

	go func() {
		f := func(text string) {
			_ = client.CurrentUser().SetStatus(&harmony.Status{
				Game: &harmony.Activity{
					Name: text,
				},
			})
		}

		if len(config.Statuses) == 1 {
			f(fmt.Sprintf(config.Statuses[0], buildInfo.Version))
			return
		}

		for {
			for _, text := range config.Statuses {
				f(fmt.Sprintf(text, buildInfo.Version))
				time.Sleep(time.Second * 15)
			}
		}
	}()

	state.AddGoroutine()

	// Set finish worker
	go func() {
		state.WaitUntilShutdownTrigger()
		client.Disconnect()
		state.FinishGoroutine()
	}()

	return nil
}