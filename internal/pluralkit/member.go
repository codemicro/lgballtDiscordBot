package pluralkit

import (
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"net/http"
)

var (
	membersBySystemIdUrl = config.PkApi.ApiUrl + "/s/%s/members"
	memberByMemberIdUrl  = config.PkApi.ApiUrl + "/m/%s"

	ErrorMemberNotFound     = errors.New("pluralkit: member with specified ID not found (PK API returned a 404)")
)

type Member struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Nickname    string `json:"display_name"`
	Description string `json:"description"`
	Pronouns    string `json:"pronouns"`
	Birthday    string `json:"birthday"`
}

func MembersBySystemId(sid string) ([]*Member, error) {
	return nil, nil
}

func MemberByMemberId(mid string) (*Member, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(memberByMemberIdUrl, mid), nil)
	if err != nil {
		return nil, err
	}

	resp, err := makeRequest(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrorMemberNotFound
	}

	mem := new(Member)
	err = parseJsonResponse(resp, mem)
	if err != nil {
		return nil, err
	}

	return mem, nil
}
