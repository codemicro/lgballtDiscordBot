package toneTags

import (
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"strings"
)

func sanitiseShorthand(tag string) string {
	// make lowercase
	tag = strings.ToLower(tag)
	// remove leading `/`
	tag = strings.TrimPrefix(tag, "/")

	return tag
}

func (*ToneTags) Lookup(ctx *route.MessageContext) error {

	tag := ctx.Arguments["tag"].(string)
	tag = sanitiseShorthand(tag)

	// lookup in DB
	dbTag := &db.ToneTag{
		Shorthand: tag,
	}

	found, err := dbTag.Get()
	if err != nil {
		return err
	}

	if !found {
		return ctx.SendErrorMessage("Tone tag shorthand not recognised")
	}

	emb := &discordgo.MessageEmbed{
		Title:       "/" + tag,
		Description: dbTag.Description,
	}

	_, err = ctx.SendMessageEmbed(ctx.Message.ChannelID, emb)
	return err
}

func (*ToneTags) List(ctx *route.MessageContext) error {

	all, err := db.GetAllToneTags()
	if err != nil {
		return err
	}

	div := len(all)/2

	var sba strings.Builder
	var sbb strings.Builder

	for i, tt := range all {
		var ts *strings.Builder
		if i < div {
			ts = &sba
		} else {
			ts = &sbb
		}
		ts.WriteString("**/")
		ts.WriteString(tt.Shorthand)
		ts.WriteString("** - ")
		ts.WriteString(tt.Description)
		ts.WriteRune('\n')
	}

	const zeroWidthSpace = "​" // yes, THERE IS SOMETHING IN THAT STRING

	emb := &discordgo.MessageEmbed{
		Title:       "Tone tag list",
		Description: "Tone tags are indicators that can be added to your message to help others in understanding your intended tone as a way to counteract the lacking nature of text-based communication. Text lacks certain audio and physical cues present in spoken speech (such as voice inflection, body language, facial expressions, etc.), so tone tags are a way to try and counteract that.\n\nHere's a list of them :P",
		Fields: []*discordgo.MessageEmbedField{
			{
				// this ZWSP is needed since Discord requires something in embed field names - I can't just leave it blank
				Name: zeroWidthSpace,
				Value:  strings.TrimSpace(sba.String()),
				Inline: true,
			},
			{
				Name: zeroWidthSpace,
				Value:  strings.TrimSpace(sbb.String()),
				Inline: true,
			},
		},
	}

	_, err = ctx.SendMessageEmbed(ctx.Message.ChannelID, emb)
	return err
}

func (*ToneTags) Create(ctx *route.MessageContext) error {

	tag := ctx.Arguments["tag"].(string)
	tag = sanitiseShorthand(tag)

	desc := ctx.Arguments["description"].(string)

	// save in DB

	dbTag := &db.ToneTag{
		Shorthand: tag,
	}
	found, _ := dbTag.Get()

	// the description is set after this call to Get() since the call will overwrite the value of dbTag.Description

	dbTag.Description = desc

	var err error
	if found {
		// already exists, update
		err = dbTag.Save()
		_ = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "📝")
	} else {
		// does not exist, create
		err = dbTag.Create()
	}

	if err != nil {
		return err
	}

	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
}

