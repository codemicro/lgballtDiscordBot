package pluralkit

import (
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/analytics"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
)

var (
	systemByIdUrl               = config.PkApi.ApiUrl + "/s/%s"
	systemByDiscordAccountIdUrl = config.PkApi.ApiUrl + "/a/%s"

	ErrorSystemNotFound     = errors.New("pluralkit: system with specified ID not found (PK API returned a 404)")
	ErrorAccountHasNoSystem = errors.New("pluralkit: account with specified ID has no systems or does not exist " +
		"(PK API returned a 404)")
)

// System represents a system object from the PluralKit API
type System struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	AvatarUrl   string `json:"avatar_url"`
}

// SystemById fetches a system from the PluralKit API based on its PluralKit ID
func SystemById(sid string) (*System, error) {
	sys := new(System)
	err := orchestrateRequest(
		fmt.Sprintf(systemByIdUrl, sid),
		sys,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorSystemNotFound},
	)
	if err != nil {
		return nil, err
	}
	analytics.ReportPluralKitRequest("System by ID")
	return sys, nil
}

// SystemByDiscordAccount fetches a system linked to a Discord user ID from the PluralKit API
func SystemByDiscordAccount(discordUid string) (*System, error) {
	sys := new(System)
	err := orchestrateRequest(
		fmt.Sprintf(systemByDiscordAccountIdUrl, discordUid),
		sys,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorAccountHasNoSystem},
	)
	if err != nil {
		return nil, err
	}
	analytics.ReportPluralKitRequest("System by Discord account ID")
	return sys, nil
}
