package misc

import (
	_ "embed"
	"encoding/json"
	"github.com/codemicro/dgo-toolkit/route"
	"math/rand"
	"regexp"
	"strings"
)

//go:embed emojifyData.json
var emojifyRawData []byte
var emojifyData map[string]map[string]int

func init() {
	err := json.Unmarshal(emojifyRawData, &emojifyData)
	if err != nil {
		panic("Cannot unmarshal emoji data: " + err.Error())
	}
}

var emojifyCommonWords = []string{
	"a",
	"an",
	"as",
	"is",
	"if",
	"of",
	"the",
	"it",
	"its",
	"or",
	"are",
	"this",
	"with",
	"so",
	"to",
	"at",
	"was",
	"and",
}

var emojifyFilterCharsRegexp = regexp.MustCompile(`[^0-9a-zA-Z]`)

func addEmojisToString(input string, density int) string {

	makeRandomChoice := func() bool { return rand.Intn(100) <= density }
	isCommonWord := func(word string) bool {
		for _, cw := range emojifyCommonWords {
			if strings.EqualFold(cw, word) {
				return true
			}
		}
		return false
	}

	words := strings.Split(strings.ReplaceAll(input, "\n", " "), " ")

	result := new(strings.Builder)

	for _, rawWord := range words {
		word := emojifyFilterCharsRegexp.ReplaceAllString(strings.ToLower(rawWord), "")

		var emojiOptions []string
		for emojis, numberOf := range emojifyData[word] {
			for i := 0; i < numberOf; i += 1 {
				emojiOptions = append(emojiOptions, emojis)
			}
		}

		result.WriteRune(' ')
		result.WriteString(rawWord)

		if isCommonWord(word) || !makeRandomChoice() || len(emojiOptions) == 0 {
			continue
		}

		var chosenEmoji string
		{
			n := rand.Intn(len(emojiOptions))
			chosenEmoji = emojiOptions[n]
		}
		result.WriteRune(' ')
		result.WriteString(chosenEmoji)
	}

	return strings.TrimSpace(result.String())
}

func (m *Misc) Emojify(ctx *route.MessageContext) error {
	// args: content
	content := ctx.Arguments["content"].(string)
	_, err := ctx.SendMessageString(ctx.Message.ChannelID, addEmojisToString(content, 80))
	return err
}