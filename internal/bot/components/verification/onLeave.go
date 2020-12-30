package verification

import (
	"context"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/skwair/harmony"
	"github.com/skwair/harmony/audit"
)

func (v *Verification) OnMemberRemove(m *harmony.GuildMemberRemove) error {
	// check audit log to see if user was kicked or banned
	baevs, err := v.b.Client.Guild(m.GuildID).AuditLog(context.Background(),
		harmony.WithEntryType(audit.EntryTypeMemberBanAdd))
	if err != nil {
		return err
	}

	kevs, err := v.b.Client.Guild(m.GuildID).AuditLog(context.Background(),
		harmony.WithEntryType(audit.EntryTypeMemberKick))
	if err != nil {
		return err
	}

	var action string
	var reason string

	// iterate backwards. most recent entries are returned first usually
	evs := append(kevs.Entries, baevs.Entries...)

	for i := len(evs) - 1; i >= 0; i -= 1 {

		switch y := evs[i].(type) {
		case *audit.MemberBanAdd:
			if y.TargetID == m.User.ID {
				action = "banned"
				reason = y.Reason
			}
		case *audit.MemberKick:
			if y.TargetID == m.User.ID && action != "banned" {
				action = "kicked"
				reason = y.Reason
			}
		}

	}

	if action == "" {
		// must just have been a user leaving
		return nil
	}

	var ur db.UserRemove
	ur.UserId = m.User.ID
	found, err := ur.Get()
	if err != nil {
		return err
	}

	ur.Action = action
	ur.Reason = reason

	if found {
		err = ur.Save()
	} else {
		err = ur.Create()
	}

	if err != nil { return err }

	// send log message
	_, err = v.b.SendMessage(modlogChannelId, "Action logged.")
	return err

}
