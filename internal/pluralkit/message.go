package pluralkit

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/analytics"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
)

var (
	messageByIdUrl = config.PkApi.ApiUrl + "/messages/%s"
)

// Message represents a message object returned from the PluralKit API
type Message struct {
	Timestamp         string  `json:"timestamp"`
	Id                string  `json:"id"`
	OriginalMessageId string  `json:"original"`
	AuthorUserId      string  `json:"sender"`
	Channel           string  `json:"channel"`
	System            *System `json:"system"`
	Member            *Member `json:"member"`
}

// MessageById fetches information about a source or proxied message from the PluralKit API
func MessageById(mid string) (*Message, error) {
	sys := new(Message)

	if err := orchestrateRequest(
		fmt.Sprintf(messageByIdUrl, mid),
		sys,
	); err != nil {
        return nil, err
    }
	// Can return MessageNotFound

	analytics.ReportPluralKitRequest("Message by ID")
	return sys, nil
}
