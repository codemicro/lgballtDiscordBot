package main

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/reddit"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"os"
	"time"
)

//go:generate python ../../scripts/prebuild.py ../.. ../../internal/buildInfo/buildInfo.go buildInfo 3.2.0

func main() {

	fmt.Printf("LGballT bot v%s built on %s (%s)\n\n", buildInfo.Version, buildInfo.BuildDate,
		buildInfo.GoVersion)

	runState := state.NewState()

	err := bot.Start(runState)
	if err != nil {
		logging.Error(err, "Failed to start Harmony client")
		os.Exit(1)
	}

	err = reddit.Start(runState)
	if err != nil {
		logging.Error(err, "Failed to start Reddit feed monitor(s)")
		os.Exit(1)
	}

	fmt.Println("Running, press ctrl+C to exit.")

	<-runState.ShutdownSignal

	fmt.Print("Shutting down... ")

	runState.TriggerShutdown()
	if runState.WaitUntilAllComplete(time.Second * 10) {
		fmt.Println()
		logging.Warn("Shutdown timeout exceeded, forcibly terminating")
	}

	fmt.Println("bye-bye!")

}
