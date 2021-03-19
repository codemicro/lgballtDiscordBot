package verification

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/meta"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"regexp"
	"time"
)

//go:generate msgp -tests=false -io=false -unexported
//msgp:ignore Verification

const (
	dataStartMarker = "[STA."
	dataEndMarker   = ".END]"

	acceptReaction = "☑️"
	rejectReaction = "❌"

	ratelimitTimeout = time.Hour
)

var (
	dataExtractionRegex = regexp.MustCompile(fmt.Sprintf("%s(.+)%s", regexp.QuoteMeta(dataStartMarker), regexp.QuoteMeta(dataEndMarker)))
	errorMissingData    = errors.New("unable to find inline data data")

	logHelpText = fmt.Sprintf("*React with %s to accept this request or %s to reject this request.\nRejecting this request will not inform the user.*", acceptReaction, rejectReaction)
)

type inlineData struct {
	UserID string `msg:"u"`
}

func dataFromString(text string) (inlineData, error) {
	matches := dataExtractionRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return inlineData{}, errorMissingData
	}

	data := matches[1]

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return inlineData{}, err
	}

	var iu inlineData
	_, err = iu.UnmarshalMsg(decoded)
	if err != nil {
		return inlineData{}, err
	}

	return iu, nil
}

func (z inlineData) toString() string {
	data, _ := z.MarshalMsg(nil)
	encoded := base64.StdEncoding.EncodeToString(data)
	return dataStartMarker + encoded + dataEndMarker
}

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
