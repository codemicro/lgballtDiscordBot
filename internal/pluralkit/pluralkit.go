package pluralkit

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
)

var (
	systemByMemberAccountIdUrl = config.PkApi.ApiUrl + "/a/%s"
	membersBySystemIdUrl = config.PkApi.ApiUrl + "/s/%s/members"
	memberByMemberIdUrl = config.PkApi.ApiUrl + "/m/%s"

	userAgent = fmt.Sprintf("r/LGBallT Discord bot v%s (%s)", buildInfo.Version, config.PkApi.ContactEmail)
)
