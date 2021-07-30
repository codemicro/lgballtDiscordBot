package pluralkit

import (
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/analytics"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
)

var (
	messageByIdUrl = config.PkApi.ApiUrl + "/msg/%s"

	ErrorMessageNotFound = errors.New("pluralkit: message with specified ID not found (PK API returned a 404)")
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
	err := orchestrateRequest(
		fmt.Sprintf(messageByIdUrl, mid),
		sys,
		func(i int) bool { return i == 200 },
		map[int]error{404: ErrorMessageNotFound},
	)
	if err != nil {
		return nil, err
	}
	analytics.ReportPluralKitRequest("Message by ID")
	return sys, nil
}
