package tools

import (
	"bytes"
	"regexp"
	"strings"
)

func GetCommand(commandString, prefix string) []string {
	command := strings.Split(strings.Trim(commandString, " "), " ")

	if len(command) >= 1 {
		if command[0] == prefix {
			command = command[1:]
		}
	}

	if len(command) == 1 && command[0] == "" {
		command = []string{}
	}
	return command
}

func IsStringInSlice(needle string, haystack []string) (found bool) {
	for _, v := range haystack {
		if strings.EqualFold(needle, v) {
			found = true
			break
		}
	}
	return
}

var messageLinkRegexp = regexp.MustCompile(`(?m)https://(.+)?discord\.com/channels/(\d+)/(\d+)/(\d+)/?`)

func ParseMessageLink(link string) (guildId, channelId, messageId string, valid bool) {

	// This is a message link: https://discord.com/channels/<guild ID>/<channel ID>/<message ID>
	matches := messageLinkRegexp.FindAllStringSubmatch(link, -1)
	if len(matches) == 0 {
		return
	}

	guildId = matches[0][2]
	channelId = matches[0][3]
	messageId = matches[0][4]
	valid = true

	return

}

var customEmojiRegex = regexp.MustCompile(`<a?:(.+:\d+)>`)

func ParseEmojiToString(eString string) string {
	// Custom emojis look like this <a:whirleythonk:743765991464501260> and match this regex: <a?:(.+:\d+)>
	var emoji string
	if customEmojiRegex.MatchString(eString) {
		matches := customEmojiRegex.FindAllStringSubmatch(eString, -1)
		emoji = matches[0][1]
	} else {
		emoji = eString
	}
	return emoji
}

var channelMentionRegex = regexp.MustCompile(`<#(\d+)>`)

func ParseChannelMention(mString string) (string, bool) {
	if channelMentionRegex.MatchString(mString) {
		matches := channelMentionRegex.FindAllStringSubmatch(mString, -1)
		return matches[0][1], true
	} else {
		return "", false
	}
}

func MakeChannelMention(channelId string) string {
	return "<#" + channelId + ">"
}

var idFromPingRegex = regexp.MustCompile(`<@!(.+)>`)

func ParsePing(ping string) (string, bool) {
	if idFromPingRegex.MatchString(ping) {
		matches := idFromPingRegex.FindAllStringSubmatch(ping, -1)
		return matches[0][1], true
	}
	return "", false
}

func MakePing(uid string) string {
	return "<@!" + uid + ">"
}

// a byte buffer that implements the Close() method
type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	// this is just memory, so all we need to do is return
	return nil
}
