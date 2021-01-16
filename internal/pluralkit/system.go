package pluralkit

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/skwair/harmony"
)

var (
	systemByIdUrl = config.PkApi.ApiUrl + "/s/%s"
	systemByMemberAccountIdUrl = config.PkApi.ApiUrl + "/a/%s"
)

type System struct {
	Id string
	Name string
	Description string
	Tag string
	AvatarUrl string
}

func SystemById(sid string) (*System, error) {
	return nil, nil
}

func SystemByDiscordAccount(m *harmony.User) (*System, error) {
	return nil, nil
}
