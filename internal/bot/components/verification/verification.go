package verification

import (
	"crypto/sha256"
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"time"
)

const (
	acceptReaction       = "✅"
	rejectReaction       = "❌"
	scrapPronounReaction = "🚮"

	pronounEmbedFieldTitle = "Pronouns"

	ratelimitTimeout = time.Hour
)

type Verification struct {
	ratelimit map[string]time.Time
}

func Init(kit *route.Kit, _ *state.State) error {

	comp := new(Verification)
	comp.ratelimit = make(map[string]time.Time)

	kit.AddCommand(&route.Command{
		Name:        "Force submit verification request",
		Help:        "Force submit a verification request for review",
		CommandText: []string{"verifyf"},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByChannel(config.VerificationIDs.InputChannel),
			route.RestrictionByRole(config.AdminRole),
		},
		Arguments: []route.Argument{
			{Name: "messageLink", Type: route.URL},
		},
		Run:       comp.FVerify,
		Invisible: true,
		Category:  meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Submit verification request",
		Help:        "Submit a verification request for review",
		CommandText: []string{"verify"},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByChannel(config.VerificationIDs.InputChannel),
		},
		Run:       comp.Verify,
		Invisible: true,
		Category:  meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Track ban",
		CommandText: []string{"ban"},
		Arguments: []route.Argument{
			{Name: "user", Type: common.PingOrUserIdType},
			{Name: "reason", Type: route.RemainingString},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRole),
		},
		Run:       comp.TrackBan,
		Invisible: true,
		Category:  meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Track kick",
		CommandText: []string{"kick"},
		Arguments: []route.Argument{
			{Name: "user", Type: common.PingOrUserIdType},
			{Name: "reason", Type: route.RemainingString},
		},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRole),
		},
		Run:       comp.TrackKick,
		Invisible: true,
		Category:  meta.CategoryAdminTools,
	})

	kit.AddCommand(&route.Command{
		Name:        "Purge unverified users",
		Help:        "Kick all users that have no assigned roles and joined more than 7 days ago",
		CommandText: []string{"purgeUnverified"},
		Restrictions: []route.CommandRestriction{
			route.RestrictionByRole(config.AdminRole),
		},
		Run:      comp.PurgeUnverifiedMembers,
		Category: meta.CategoryAdminTools,
	})

	kit.AddReaction(&route.Reaction{
		Name:  "Verification reaction handler",
		Run:   comp.DecisionReaction,
		Event: route.ReactionAdd,
	})

	return nil
}

func hashString(raw string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(raw)))
}
