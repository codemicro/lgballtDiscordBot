package actionLog

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/codemicro/lgballtDiscordBot/internal/pluralkit"
	"github.com/codemicro/lgballtDiscordBot/internal/tools"
	"strings"
	"time"
)

func messageUpdate(s *discordgo.Session, mu *discordgo.MessageUpdate) {

	if wasBot(mu.Message) {
		return
	}

	// ignore PK proxied things
	pkMessage, err := pluralkit.MessageById(mu.ID)
	if err != nil {
		// err != nil
		if !errors.Is(err, pluralkit.ErrorMessageNotFound) {
			logging.Error(err, "messageUpdate actionLog handler")
		}
	}

	var files []*discordgo.File
	var sb strings.Builder

	sb.WriteString("**Message updated** from ")
	sb.WriteString(getAuthorMention(mu.Message))

	if pkMessage != nil {
		sb.WriteString(" (")
		sb.WriteString(tools.MakePing(pkMessage.AuthorUserId))
		sb.WriteString(")")
	}

	sb.WriteString(" in ")
	sb.WriteString(tools.MakeChannelMention(mu.ChannelID))

	previousContent := getContent(mu.BeforeUpdate)
	newContent := getContent(mu.Message)

	if len(previousContent)+len(newContent) > 1500 {
		files = append(files, &discordgo.File{
			Name:        "before.txt",
			ContentType: "text/plain",
			Reader:      strings.NewReader(previousContent),
		}, &discordgo.File{
			Name:        "after.txt",
			ContentType: "text/plain",
			Reader:      strings.NewReader(newContent),
		})
	} else {
		sb.WriteString("\n**Old**: `")
		sb.WriteString(previousContent)
		sb.WriteString("`\n**New**: `")
		sb.WriteString(newContent)
		sb.WriteString("`\n")
	}

	if err = log(s, eventTypeMessageUpdate, sb.String(), files...); err != nil {
		logging.Error(err, "messageDelete actionLog handler")
	}

}

func messageDelete(s *discordgo.Session, md *discordgo.MessageDelete) {

	if wasBot(md.Message) {
		return
	}

	// ignore verification
	if md.ChannelID == config.VerificationIDs.InputChannel {
		return
	}

	// ignore PK proxied things
	pkMessage, err := pluralkit.MessageById(md.ID)
	if err == nil {
		if pkMessage.OriginalMessageId == md.ID {
			// the pluralkit proxy message has been deleted
			return
		}
	} else {
		// err != nil
		if !errors.Is(err, pluralkit.ErrorMessageNotFound) {
			logging.Error(err, "messageDelete actionLog handler")
		}
	}

	var files []*discordgo.File
	var sb strings.Builder

	sb.WriteString("**Message deleted** from ")
	sb.WriteString(getAuthorMention(md.BeforeDelete))

	if pkMessage != nil {
		sb.WriteString(" (")
		sb.WriteString(tools.MakePing(pkMessage.AuthorUserId))
		sb.WriteString(")")
	}

	sb.WriteString(" in ")
	sb.WriteString(tools.MakeChannelMention(md.ChannelID))

	messageContent := getContent(md.BeforeDelete)

	if len(messageContent) > 1500 {
		files = append(files, &discordgo.File{
			Name:        "message.txt",
			ContentType: "text/plain",
			Reader:      strings.NewReader(messageContent),
		})
	} else {
		sb.WriteString(": `")
		sb.WriteString(messageContent)
		sb.WriteString("`")
	}

	err = log(s, eventTypeMessageDelete, sb.String(), files...)
	if err != nil {
		logging.Error(err, "messageDelete actionLog handler")
	}
}

func messageDeleteBulk(s *discordgo.Session, mdb *discordgo.MessageDeleteBulk) {
	var messageBuilder strings.Builder

	messageBuilder.WriteString("**Multiple messages deleted** in ")
	messageBuilder.WriteString(tools.MakeChannelMention(mdb.ChannelID))

	var fileBuilder strings.Builder

	for i := len(mdb.Messages) - 1; i >= 0; i -= 1 {
		messageID := mdb.Messages[i]

		fileBuilder.WriteString(messageID)
		fileBuilder.WriteRune(' ')

		timestamp, _ := discordgo.SnowflakeTimestamp(messageID)
		fileBuilder.WriteString(timestamp.Format(time.RFC822))

		fromState := mdb.BeforeDelete[messageID]

		fileBuilder.WriteRune(' ')

		fileBuilder.WriteString(getAuthorUsername(fromState))
		fileBuilder.WriteString(": ")
		fileBuilder.WriteString(getContent(fromState))
		fileBuilder.WriteRune('\n')
	}

	err := log(s, eventTypeMessageDelete, messageBuilder.String(), &discordgo.File{
		Name:        "messages.txt",
		ContentType: "text/plain",
		Reader:      strings.NewReader(fileBuilder.String()),
	})

	if err != nil {
		logging.Error(err, "messageDelete actionLog handler")
	}
}
