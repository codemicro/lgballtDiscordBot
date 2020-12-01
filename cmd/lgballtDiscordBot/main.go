package main

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/bios"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/info"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/skwair/harmony"
	"os"
	"os/signal"
	"time"
)

const (
	version = "1.0.1"
)

func main() {

	fmt.Printf("LGballT bot v%s\n\n", version)

	client, err := harmony.NewClient(config.Config.Token)
	if err != nil {
		fmt.Printf("Failed to initialise a new Harmony client\n\n")
		os.Exit(1)
	}

	b := bot.New(client, config.Config.Prefix)

	if err := bios.Register(b, "bio"); err != nil {
		logging.Error(err, "Failed to register component bio")
		os.Exit(1)
	}

	if err := info.Register(b, "info"); err != nil {
		logging.Error(err, "Failed to register component info")
		os.Exit(1)
	}

	if err = client.Connect(context.Background()); err != nil {
		logging.Error(err, "Failed to connect to Discord")
		os.Exit(1)
	}
	defer client.Disconnect()

	go func() {
		f := func(text string) {
			_ = client.CurrentUser().SetStatus(&harmony.Status{
				Game: &harmony.Activity{
					Name: text,
				},
			})
		}

		if len(config.Config.Statuses) == 1 {
			f(fmt.Sprintf(config.Config.Statuses[0], version))
			return
		}

		for {
			for _, text := range config.Config.Statuses {
				f(fmt.Sprintf(text, version))
				time.Sleep(time.Second * 15)
			}
		}
	}()


	fmt.Println("Bot is running, press ctrl+C to exit.")

	// Wait for ctrl-C, then exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Println("Waiting 5 seconds to ensure data is saved...")

	time.Sleep(time.Second * 5)

	fmt.Println("Shutting down - bye-bye!")

}