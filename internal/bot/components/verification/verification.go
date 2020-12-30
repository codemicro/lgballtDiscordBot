package verification

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/core"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"regexp"
)

//go:generate msgp -tests=false -io=false -unexported
//msgp:ignore Verification

const (
	dataStartMarker = "[STA."
	dataEndMarker   = ".END]"

	acceptReaction = "☑️"
	rejectReaction = "❌"
)

var (
	roleId          string
	InputChannelId  string
	OutputChannelId string
	modlogChannelId string
)

func init() {
	roleId = config.VerificationIDs.RoleId
	InputChannelId = config.VerificationIDs.InputChannel
	OutputChannelId = config.VerificationIDs.OutputChannel
	modlogChannelId = config.VerificationIDs.ModlogChannel
}

var (
	dataExtractionRegex = regexp.MustCompile(fmt.Sprintf("%s(.+)%s", regexp.QuoteMeta(dataStartMarker), regexp.QuoteMeta(dataEndMarker)))
	errorMissingData    = errors.New("unable to find inline data data")

	logHelpText = fmt.Sprintf("*React with %s to accept this request or %s to reject this request.\nRejecting this request will not inform the user.*", acceptReaction, rejectReaction)
)

type Verification struct {
	b *core.Bot
}

func New(bot *core.Bot) (*Verification, error) {
	b := new(Verification)
	b.b = bot

	return b, nil
}

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
