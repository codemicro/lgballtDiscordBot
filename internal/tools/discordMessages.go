package tools

import (
	"fmt"
	"regexp"
	"strings"
)

var messageLinkRegexp = regexp.MustCompile(`(?m)https://(.+)?discord(?:app)?\.com/channels/(\d+)/(\d+)/(\d+)/?`)

func ParseMessageLink(link string) (guildId, channelId, messageId string, valid bool) {

	// This is a message link: https://discord.com/channels/<guild ID>/<channel ID>/<message ID>
	matches := messageLinkRegexp.FindStringSubmatch(link)
	if len(matches) == 0 {
		return
	}

	guildId = matches[2]
	channelId = matches[3]
	messageId = matches[4]
	valid = true

	return

}

func MakeMessageLink(guildId, channelId, messageId string) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildId, channelId, messageId)
}

var CustomEmojiRegex = regexp.MustCompile(`<(a?):([^:]+:\d+)>`)

// ParseEmojiToString returns the emoji name and ID in a single string like name:id
func ParseEmojiToString(eString string) string {
	// Custom emojis look like this <a:whirleythonk:743765991464501260> and match this regex: <a?:(.+:\d+)>
	var emoji string
	if CustomEmojiRegex.MatchString(eString) {
		matches := CustomEmojiRegex.FindAllStringSubmatch(eString, -1)
		emoji = matches[0][2]
	} else {
		emoji = eString
	}
	return emoji
}

// ParseEmojiComponents returns if the emoji is animated, the emoji name and emoji ID
func ParseEmojiComponents(eString string) (isValid bool, isAnimated bool, emojiName string, emojiID string) {
	if CustomEmojiRegex.MatchString(eString) {
		matches := CustomEmojiRegex.FindAllStringSubmatch(eString, -1)
		aStr := matches[0][1]
		combi := strings.Split(matches[0][2], ":")
		if aStr == "a" {
			isAnimated = true
		}
		return true, isAnimated, combi[0], combi[1]
	}
	return false, false, eString, ""
}

func MakeCustomEmoji(animated bool, name, id string) string {
	o := "<"
	if animated {
		o += "a"
	}
	o += fmt.Sprintf(":%s:%s>", name, id)
	return o
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

var idFromPingRegex = regexp.MustCompile(`<@!?(.+)>`)

func ParsePing(ping string) (string, bool) {
	if idFromPingRegex.MatchString(ping) {
		matches := idFromPingRegex.FindStringSubmatch(ping)
		return matches[1], true
	}
	return ping, false
}

func MakePing(uid string) string {
	return "<@" + uid + ">"
}

func MakeRolePing(rid string) string {
	return "<@&" + rid + ">"
}

var rolePingRegex = regexp.MustCompile(`<@&(\d+)>`)

func FindRolePings(in string) (o []string) {
	for _, x := range rolePingRegex.FindAllStringSubmatch(in, -1) {
		o = append(o, x[1])
	}
	return
}

func FilterRolePing(instr string) string {
	return rolePingRegex.ReplaceAllString(instr, "`$0`")
}
