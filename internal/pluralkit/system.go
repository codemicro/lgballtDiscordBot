package pluralkit

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/analytics"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
)

var (
	systemByIdUrl = config.PkApi.ApiUrl + "/systems/%s"
)

// System represents a system object from the PluralKit API
type System struct {
	UUID        string `json:"uuid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	AvatarUrl   string `json:"avatar_url"`
}

// SystemById fetches a system from the PluralKit API based on its PluralKit ID
func SystemById(sid string) (*System, error) {
	sys := new(System)

	if err := orchestrateRequest(
		fmt.Sprintf(systemByIdUrl, sid),
		sys,
		); err != nil {
		return nil, err
	}
	// Can return SystemNotFound

	analytics.ReportPluralKitRequest("System by ID")
	return sys, nil
}

