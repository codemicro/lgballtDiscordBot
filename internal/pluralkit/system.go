package pluralkit

import (
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/skwair/harmony"
	"net/http"
)

var (
	systemByIdUrl               = config.PkApi.ApiUrl + "/s/%s"
	systemByDiscordAccountIdUrl = config.PkApi.ApiUrl + "/a/%s"

	ErrorSystemNotFound     = errors.New("pluralkit: system with specified ID not found (PK API returned a 404)")
	ErrorAccountHasNoSystem = errors.New("pluralkit: account with specified ID has no systems or does not exist " +
		"(PK API returned a 404)")
)

type System struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	AvatarUrl   string `json:"avatar_url"`
}

func SystemById(sid string) (*System, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(systemByIdUrl, sid), nil)
	if err != nil {
		return nil, err
	}

	resp, err := makeRequest(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrorSystemNotFound
	}

	sys := new(System)
	err = parseJsonResponse(resp, sys)
	if err != nil {
		return nil, err
	}

	return sys, nil
}

func SystemByDiscordAccount(m *harmony.User) (*System, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(systemByDiscordAccountIdUrl, m.ID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := makeRequest(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrorAccountHasNoSystem
	}

	sys := new(System)
	err = parseJsonResponse(resp, sys)
	if err != nil {
		return nil, err
	}

	return sys, nil
}
