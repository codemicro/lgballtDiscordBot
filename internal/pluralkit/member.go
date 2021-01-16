package pluralkit

import "github.com/codemicro/lgballtDiscordBot/internal/config"

var (
	membersBySystemIdUrl = config.PkApi.ApiUrl + "/s/%s/members"
	memberByMemberIdUrl = config.PkApi.ApiUrl + "/m/%s"
)

type Member struct {
	Id string
	Name string
	Nickname string
	Description string
	Pronouns string
	Birthday string
}

func MembersBySystemId(sid string) ([]*Member, error) {
	return nil, nil
}

func MemberByMemberId(mid string) (*Member, error) {
	return nil, nil
}
