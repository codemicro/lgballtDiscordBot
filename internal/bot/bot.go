package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/bios"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/chatchart"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/info"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/messageTools"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/misc"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/muteme"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/pressf"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/roles"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"time"
)

func Start(state *state.State) error {

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return err
	}

	kit := route.NewKit(session, []string{config.Prefix})

	kit.ErrorHandler = func(err error) { logging.Error(err) }

	err = registerHandlers(kit, state)
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}

	go func() {
		f := func(text string) {
			err := session.UpdateGameStatus(0, text)
			if err != nil {
				logging.Error(err)
			}
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

	go func() {
		state.WaitUntilShutdownTrigger()
		_ = session.Close()
		state.FinishGoroutine()
	}()

	return nil
}

func registerHandlers(kit *route.Kit, st *state.State) error {

	// TODO: commands/reactions

	err := misc.Init(kit, st)
	if err != nil {
		return err
	}

	err = pressf.Init(kit, st)
	if err != nil {
		return err
	}

	err = muteme.Init(kit, st)
	if err != nil {
		return err
	}

	err = info.Init(kit, st)
	if err != nil {
		return err
	}

	err = bios.Init(kit, st)
	if err != nil {
		return err
	}

	err = messageTools.Init(kit, st)
	if err != nil {
		return err
	}

	err = chatchart.Init(kit, st)
	if err != nil {
		return err
	}

	err = roles.Init(kit, st)
	if err != nil {
		return err
	}

	kit.CreateHandlers()

	return nil
}
