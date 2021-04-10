package verification

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"time"
)

//go:generate msgp -tests=false -io=false -unexported
//msgp:ignore Verification

const (
	dataStartMarker = "[STA."
	dataEndMarker   = ".END]"

	acceptReaction = "✅"
	rejectReaction = "❌"
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
		Category: meta.CategoryAdminTools,
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
		Category: meta.CategoryAdminTools,
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
		Category: meta.CategoryAdminTools,
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
		Category: meta.CategoryAdminTools,
	})

	kit.AddReaction(&route.Reaction{
		Name:  "Verification reaction handler",
		Run:   comp.DecisionReaction,
		Event: route.ReactionAdd,
	})

	return nil
}
