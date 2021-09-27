package misc

import (
	"bytes"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"io"
	"io/ioutil"
	"net/http"
)

const maxUploadSize = 8 * 1000 * 1000 // 8MB

func (*Misc) SpoilerThis(ctx *route.MessageContext) error {

	if len(ctx.Message.Attachments) == 0 {
		return ctx.SendErrorMessage("there must be at least one attachment")
	}

	var nf []*discordgo.File
	var totalFileSize int

	for _, attachment := range ctx.Message.Attachments {

		if totalFileSize+attachment.Size > maxUploadSize {
			return ctx.SendErrorMessage("attachments are too large and cannot be re-uploaded (greater than 7MB)")
		}

		fileBody, err := downloadFile(attachment.URL)
		if err != nil {
			return nil
		}

		firstSection := make([]byte, 512)
		_, err = fileBody.Read(firstSection)
		if err != nil {
			return err
		}

		remaining, err := ioutil.ReadAll(fileBody)
		if err != nil {
			return err
		}

		_ = fileBody.Close()

		nf = append(nf, &discordgo.File{
			Name:        "SPOILER_" + attachment.Filename,
			ContentType: http.DetectContentType(firstSection),
			Reader:      bytes.NewReader(append(firstSection, remaining...)),
		})
		totalFileSize += attachment.Size

	}

	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Content:         "From " + ctx.Message.Author.Mention(),
		Files:           nf,
		AllowedMentions: ctx.DefaultAllowedMentions(),
	})

	if err != nil {
		return err
	}

	return ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
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
