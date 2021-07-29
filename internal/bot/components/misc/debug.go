package misc

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"runtime/pprof"
	"time"
)

func (m *Misc) GoroutineStack(ctx *route.MessageContext) error {

	b := new(bytes.Buffer)
	_ = pprof.Lookup("goroutine").WriteTo(b, 1)

	// make DM channel
	dm, err := ctx.Session.UserChannelCreate(ctx.Message.Author.ID)
	if err != nil {
		return err
	}

	_, err = ctx.Session.ChannelMessageSendComplex(dm.ID, &discordgo.MessageSend{
		Content: time.Now().Format(time.RFC1123),
		Files: []*discordgo.File{
			{
				Name:        fmt.Sprintf("goroutineStacktrace-%d-.txt", time.Now().Unix()),
				ContentType: "text/plain",
				Reader:      b,
			},
		},
	})

	return err
}
