package main

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/bios"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/skwair/harmony"
	"io/ioutil"
	"os"
	"os/signal"
	"time"
)

const (
	version = "1.0.0"
	botPrefix = "*"
)

func main() {

	fmt.Printf("LGballT bot v%s\n\n", version)

	tokenFileName := "token.txt"
	tokenBytes, err := ioutil.ReadFile(tokenFileName)
	if err != nil {
		logging.Error(err, fmt.Sprintf("Failed to read %s", tokenFileName))
		os.Exit(1)
	}

	client, err := harmony.NewClient(string(tokenBytes))
	if err != nil {
		fmt.Printf("Failed to initialise a new Harmony client\n\n")
		os.Exit(1)
	}

	b := bot.New(client, botPrefix)

	if err := bios.Register(b, "bio"); err != nil {
		logging.Error(err, "Failed to register component bio")
		os.Exit(1)
	}

	if err = client.Connect(context.Background()); err != nil {
		logging.Error(err, "Failed to connect to Discord")
		os.Exit(1)
	}
	defer client.Disconnect()

	_ = client.CurrentUser().SetStatus(&harmony.Status{
		Game: &harmony.Activity{
			Name: fmt.Sprintf("v%s - ping 0x5444#8669 if something breaks", version),
		},
	})

	fmt.Println("Bot is running, press ctrl+C to exit.")

	// Wait for ctrl-C, then exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Println("Waiting 5 seconds to ensure data is saved...")

	time.Sleep(time.Second * 5)

	fmt.Println("Shutting down - bye-bye!")

}