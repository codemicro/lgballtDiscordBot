package misc

import (
	"context"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/embed"
	"sync"
	"time"
)

var (
	// channel ID, user ID
	activeMuteRequests = make(map[string]muteMeRequest)
	amrMux             sync.RWMutex

	muteMeText = "This will mute you until %s and cannot be undone. Are you sure?"
)

type muteMeRequest struct {
	userId  string
	endTime time.Time
}

func (s *Misc) startMuteRemovalWorker() {

	s.b.State.AddGoroutine()

	ticker := time.NewTicker(time.Minute)
	finished := make(chan bool)

	go func() {
		s.b.State.WaitUntilShutdownTrigger()
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
			s.muteRemovalWorker()
		}
	}

	s.b.State.FinishGoroutine()
}

func (s *Misc) muteRemovalWorker() {
	mutes, err := db.GetAllUserMutes()
	if err != nil {
		logging.Error(err, "misc.Misc.muteRemovalWorker: failed to fetch all user mutes")
		return
	}

	for _, mute := range mutes {
		if time.Now().Unix() > mute.EndTime {
			// mute expired, remove mute
			// get guild
			guild := s.b.Client.Guild(mute.GuildId)
			// remove timeout role
			err = guild.RemoveMemberRole(context.Background(), mute.UserId, config.MuteMe.TimeoutRole)
			if err != nil {
				logging.Error(err, "misc.Misc.muteRemovalWorker: failed to remove timeout role")
				continue
			}

			// adding removed roles
			for _, role := range mute.RemovedRoles {
				err = guild.AddMemberRole(context.Background(), mute.UserId, role)
				if err != nil {
					logging.Error(err, "misc.Misc.muteRemovalWorker: failed to add removed user role")
					continue
				}
			}

			// delete from database
			err = mute.Delete()
			if err != nil {
				logging.Error(err, "misc.Misc.muteRemovalWorker: failed to delete UserMute entry from DB")
				continue
			}
		}
	}
}

func (s *Misc) MuteMe(args []string, m *harmony.Message) error {

	if len(args) < 1 {
		_, err := s.b.SendMessage(m.ChannelID, "Missing duration")
		return err
	}

	duration := tools.ParseDuration(args[0])
	finishTime := time.Now().Add(duration)
	finishTimeString := finishTime.UTC().Format(time.RFC822)

	if finishTime.After(time.Now().Add(time.Hour * 24 * 366)) {
		_, err := s.b.SendMessage(m.ChannelID, "Maximum duration is 365 days")
		return err
	}

	emb := embed.Embed{
		Type:        "rich",
		Title:       "Confirmation",
		Footer:      embed.NewFooter().Text(fmt.Sprintf("React to this message with %s if you wish to mute yourself. (%s to cancel)", acceptReaction, rejectReaction)).Build(),
		Description: fmt.Sprintf(muteMeText, finishTimeString),
	}

	msg, err := s.b.SendEmbed(m.ChannelID, &emb)
	if err != nil {
		return err
	}

	amrMux.Lock()
	activeMuteRequests[msg.ID] = muteMeRequest{
		userId:  m.Author.ID,
		endTime: finishTime,
	}
	amrMux.Unlock()

	for _, v := range []string{acceptReaction, rejectReaction} {
		err := s.b.Client.Channel(msg.ChannelID).AddReaction(context.Background(), msg.ID, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Misc) MuteMeReaction(r *harmony.MessageReaction) error {

	if !(r.Emoji.Name == acceptReaction || r.Emoji.Name == rejectReaction) {
		return nil
	}

	amrMux.RLock()
	muteRequest, found := activeMuteRequests[r.MessageID]
	amrMux.RUnlock()

	if !found {
		return nil
	}

	if r.UserID != muteRequest.userId {
		return nil
	}

	err := s.b.Client.Channel(r.ChannelID).RemoveAllReactions(context.Background(), r.MessageID)
	if err != nil {
		return err
	}

	alrMux.Lock()
	delete(activeMuteRequests, r.MessageID)
	alrMux.Unlock()

	if r.Emoji.Name == acceptReaction {
		// do mute

		// get guild
		guild := s.b.Client.Guild(r.GuildID)

		// get user
		user, err := guild.Member(context.Background(), muteRequest.userId)
		if err != nil {
			return err
		}
		// determine which roles need to be removed
		var rolesToRemove []string
		for _, role := range config.MuteMe.RolesToRemove {
			if tools.IsStringInSlice(role, user.Roles) {
				rolesToRemove = append(rolesToRemove, role)
			}
		}

		// add DB entry
		um := db.UserMute{
			UserId:       muteRequest.userId,
			GuildId:      r.GuildID,
			EndTime:      muteRequest.endTime.Unix(),
			RemovedRoles: rolesToRemove,
		}
		err = um.Create()
		if err != nil {
			return err
		}

		// add timeout role
		err = guild.AddMemberRole(context.Background(), muteRequest.userId, config.MuteMe.TimeoutRole)
		if err != nil {
			return err
		}

		// remove roles
		for _, role := range rolesToRemove {
			err = guild.RemoveMemberRole(context.Background(), muteRequest.userId, role)
			if err != nil {
				return err
			}
		}

		_, _ = s.b.Client.Channel(r.ChannelID).EditMessage(context.Background(), r.MessageID, "Muted!")
	} else {
		_, _ = s.b.Client.Channel(r.ChannelID).EditMessage(context.Background(), r.MessageID, "Cancelled!")
	}

	return nil
}
