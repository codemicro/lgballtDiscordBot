package common

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"strings"
)

func isStringInSlice(needle string, haystack []string) (found bool) {
	for _, v := range haystack {
		if strings.EqualFold(needle, v) {
			found = true
			break
		}
	}
	return
}

func IsAdmin(userRoles []string) bool {
	for _, roleId := range config.AdminRoles {
		if isStringInSlice(roleId, userRoles) {
			return true
		}
	}
	return false
}

func IsOwner(userID string) bool {
	return isStringInSlice(userID, config.OwnerIds)
}