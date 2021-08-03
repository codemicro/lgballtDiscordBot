package muteme

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/rs/zerolog/log"
	"time"
)

func (mm *MuteMe) startMuteRemovalWorker(session *discordgo.Session, st *state.State) {

	ticker := time.NewTicker(time.Minute)
	finished := make(chan bool)

	go func() {
		st.WaitUntilShutdownTrigger()
		ticker.Stop()
		finished <- true
	}()

	var jumpOut bool
	for {

		if jumpOut {
			break
		}

		// So... apparently it's possible to break out of a select statement
		// Half an hour wasted debugging shutdown deadlocks... :)
		select {
		case <-finished:
			jumpOut = true
		case <-ticker.C:
			mm.muteRemovalWorker(session)
		}
	}

	st.FinishGoroutine()
}

func (mm *MuteMe) muteRemovalWorker(session *discordgo.Session) {
	
	logger := log.With().Str("area", "muteRemovalWorker").Logger()
	
	mutes, err := db.GetAllUserMutes()
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch all user mutes")
		return
	}

	for _, mute := range mutes {
		if time.Now().Unix() > mute.EndTime {
			// mute expired, remove mute
			// remove timeout role
			err = session.GuildMemberRoleRemove(mute.GuildId, mute.UserId, config.MuteMe.TimeoutRole)
			if err != nil {
				logger.Error().Err(err).Msg("failed to remove timeout role")
				continue
			}

			// adding removed roles
			for _, role := range mute.RemovedRoles {
				err = session.GuildMemberRoleAdd(mute.GuildId, mute.UserId, role)
				if err != nil {
					logger.Error().Err(err).Msg("failed to add removed user role")
					continue
				}
			}

			// delete from database
			err = mute.Delete()
			if err != nil {
				logger.Error().Err(err).Msg("failed to delete UserMute entry from DB")
				continue
			}
		}
	}
}
