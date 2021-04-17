package pluralkit

import (
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
)

var (
	membersBySystemIdUrl = config.PkApi.ApiUrl + "/s/%s/members"
	memberByMemberIdUrl  = config.PkApi.ApiUrl + "/m/%s"

	ErrorMemberNotFound    = errors.New("pluralkit: member with specified ID not found (PK API returned a 404)")
	ErrorMemberListPrivate = errors.New("pluralkit: target system found but member list is private (PK API " +
		"returned a 403)")
)

type Member struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Nickname    string `json:"display_name"`
	Avatar      string `json:"avatar_url"`
	Description string `json:"description"`
	Pronouns    string `json:"pronouns"`
	Birthday    string `json:"birthday"`
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
	return members, orchestrateRequest(
		fmt.Sprintf(membersBySystemIdUrl, sid),
		&members,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorSystemNotFound, 403: ErrorMemberListPrivate},
	)
}

func MemberByMemberId(mid string) (*Member, error) {
	member := new(Member)
	return member, orchestrateRequest(
		fmt.Sprintf(memberByMemberIdUrl, mid),
		member,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorMemberNotFound},
	)
}
