package templates

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"net/url"
)

func ViewURL(bio db.UserBio) string {
	viewURL := fmt.Sprintf("/bio/view?user=%s", url.PathEscape(bio.UserId))
	if bio.SysMemberID != "" {
		viewURL += "&member=" + url.PathEscape(bio.SysMemberID)
	}
	return viewURL
}
