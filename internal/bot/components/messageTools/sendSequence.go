package messageTools

import (
	"errors"
	"fmt"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/markdown"
	"io"
	"io/ioutil"
	"net/http"
)

func (m *MessageTools) SendSequence(ctx *route.MessageContext) error {

	channelID := ctx.Arguments["channelID"].(string)

	if len(ctx.Message.Attachments) != 1 {
		return ctx.SendErrorMessage("there must be exactly one attachment")
	}

	file := ctx.Message.Attachments[0]

	fileBody, err := downloadFile(file.URL)
	if err != nil {
		return nil
	}

	contents, err := ioutil.ReadAll(fileBody)
	if err != nil {
		return err
	}

	_ = fileBody.Close()

	markdownSections := markdown.SplitByHeader(string(contents))

	for _, section := range markdownSections {
		_, err := ctx.SendMessageString(channelID, fmt.Sprintf("**%s**\n%s", section.Title, section.Content))
		if err != nil {
			return err
		}
	}

	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
}

func downloadFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("non-200 status code returned")
	}
	return resp.Body, nil
}