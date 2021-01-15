package main

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/reddit"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"os"
	"os/signal"
	"time"
)

//go:generate python ../../scripts/prebuild.py ../.. ../../internal/buildInfo/buildInfo.go buildInfo 2.0.1

func main() {

	fmt.Printf("LGballT bot v%s built on %s (%s)\n\n", buildInfo.Version, buildInfo.BuildDate,
		buildInfo.GoVersion)

	state := tools.NewState()

	err := bot.Start(state)
	if err != nil {
		logging.Error(err, "Failed to start Harmony client")
		os.Exit(1)
	}

	err = reddit.Start(state)
	if err != nil {
		logging.Error(err, "Failed to start Reddit feed monitor(s)")
		os.Exit(1)
	}

	fmt.Println("Running, press ctrl+C to exit.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Print("Shutting down... ")

	state.TriggerShutdown()
	if state.WaitUntilAllComplete(time.Second * 10) {
		fmt.Println()
		logging.Warn("Shutdown timeout exceeded, forcibly terminating")
	}

	fmt.Println("bye-bye!")

}
