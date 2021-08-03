package actionLog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

func Setup(s *discordgo.Session) error {

	logger = log.With().Str("area", "actionLog").Logger()

	s.AddHandler(messageUpdate)
	s.AddHandler(messageDelete)
	s.AddHandler(messageDeleteBulk)

	return nil
}

const (
	eventTypeMessageDelete = iota
	eventTypeMessageUpdate
)

func logEvent(s *discordgo.Session, eventType uint8, text string, files ...*discordgo.File) error {

	var emoji string

	switch eventType {
	case eventTypeMessageDelete:
		emoji = "🗑️"
	case eventTypeMessageUpdate:
		emoji = "📝"
	}

	_, err := s.ChannelMessageSendComplex(config.ActionLogChannel, &discordgo.MessageSend{
		Content: fmt.Sprintf("%s %s", emoji, text),
		AllowedMentions: &discordgo.MessageAllowedMentions{},
		Files: files,
	})

	return err
}
