package pluralkit

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/analytics"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
)

var (
	membersBySystemIdUrl = config.PkApi.ApiUrl + "/systems/%s/members"
	memberByMemberIdUrl  = config.PkApi.ApiUrl + "/members/%s"
)

type Member struct {
	UUID        string `json:"uuid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Nickname    string `json:"display_name"`
	Avatar      string `json:"avatar_url"`
	Description string `json:"description"`
	Pronouns    string `json:"pronouns"`
	Birthday    string `json:"birthday"`
	Colour      string `json:"color"`
}

type Members []*Member

func (m Members) Get(memberID string) *Member {
	for _, member := range m {
		if member.Id == memberID {
			return member
		}
	}
	return nil
}

func MembersBySystemId(sid string) (Members, error) {
	var members Members

	if err := orchestrateRequest(
		fmt.Sprintf(membersBySystemIdUrl, sid),
		&members,
	); err != nil {
		return nil, err
	}
	// Can return SystemNotFound or MemberListPrivate

	analytics.ReportPluralKitRequest("Members by system ID")
	return members, nil
}

func MemberByMemberId(mid string) (*Member, error) {
	member := new(Member)
	if err := orchestrateRequest(
		fmt.Sprintf(memberByMemberIdUrl, mid),
		&member,
	); err != nil {
		return nil, err
	}
	// Can return MemberNotFound

	analytics.ReportPluralKitRequest("Member by ID")
	return member, nil
}
