package main

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/adminSite"
	"github.com/codemicro/lgballtDiscordBot/internal/analytics"
	"github.com/codemicro/lgballtDiscordBot/internal/bot"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/reddit"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func main() {

	fmt.Printf("LGballT bot v%s built on %s (%s)\n\n", buildInfo.Version, buildInfo.BuildDate,
		buildInfo.GoVersion)

	setupLogging(config.DebugMode, "lgballtBot.log")

	runState := state.NewState()

	err := bot.Start(runState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start bot")
		os.Exit(1)
	}

	err = reddit.Start(runState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start Reddit feed monitor(s)")
		os.Exit(1)
	}

	err = analytics.Start(runState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start analytics service")
		os.Exit(1)
	}

	err = adminSite.Start(runState)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start admin site service")
		os.Exit(1)
	}

	fmt.Println("Running, press ctrl+C to exit.")

	<-runState.ShutdownSignal

	fmt.Print("Shutting down... ")

	runState.TriggerShutdown()
	if runState.WaitUntilAllComplete(time.Second * 10) {
		fmt.Println()
		log.Warn().Msg("Shutdown timeout exceeded, forcibly terminating")
	}

	fmt.Println("bye-bye!")

}

func setupLogging(debug bool, logFile string) {

	var writer io.Writer
	var logLevel zerolog.Level

	lumberjackWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}

	if debug {
		writer = io.MultiWriter(lumberjackWriter, zerolog.NewConsoleWriter())
		logLevel = zerolog.DebugLevel
	} else {
		writer = lumberjackWriter
		logLevel = zerolog.InfoLevel
	}

	log.Logger = zerolog.New(writer).Level(logLevel).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		unpacked := eris.Unpack(err)
		var out []map[string]interface{}
		for _, chainItem := range unpacked.ErrChain {
			x := map[string]interface{}{
				"func": chainItem.Frame.Name,
				"file": chainItem.Frame.File,
				"line": chainItem.Frame.Line,
			}

			if chainItem.Msg != "" {
				x["msg"] = chainItem.Msg
			}

			out = append(out, x)
		}
		return out
	}
}