package toneTags

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"strings"
	"time"
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

	const (
		zeroWidthSpace   = "​" // yes, THERE IS SOMETHING IN THAT STRING
		embedTitle       = "Tone tag list"
		embedDescription = "Tone tags are indicators that can be added to your message to help others in understanding your intended tone as a way to counteract the lacking nature of text-based communication. Text lacks certain audio and physical cues present in spoken speech (such as voice inflection, body language, facial expressions, etc.), so tone tags are a way to try and counteract that.\n\nHere's a list of them :P"
	)

	var (
		embeds            []*discordgo.MessageEmbed
		currentEmbed      *discordgo.MessageEmbed
		currentFieldIndex int

		newEmbed = func() {
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       embedTitle,
				Description: embedDescription,
				Fields:      []*discordgo.MessageEmbedField{},
			})
			currentEmbed = embeds[len(embeds)-1]
		}

		newField = func() {
			currentEmbed.Fields = append(currentEmbed.Fields, &discordgo.MessageEmbedField{
				Name:   zeroWidthSpace,
				Value:  "",
				Inline: true,
			})
			currentFieldIndex = len(currentEmbed.Fields) - 1
		}
	)

	newEmbed()
	newField()

	for _, tt := range all {
		line := fmt.Sprintf("**/%s** - %s", tt.Shorthand, tt.Description)

		// if the current field is going to be full (when the field length hits 1024, it trips up Discord, but here we
		// have the limit set a little bit lower to make the embed more readable)
		if len(currentEmbed.Fields[currentFieldIndex].Value)+len(line)+1 >= 512 { // plus one for the \n
			// if no more fields should be added, so we make a new embed then add fields to that
			if len(currentEmbed.Fields) >= 2 {
				newEmbed()
			}

			newField()
		}

		currentEmbed.Fields[currentFieldIndex].Value += "\n" + line
	}

	if len(embeds) == 1 {
		_, err = ctx.SendMessageEmbed(ctx.Message.ChannelID, embeds[0])
		return err
	}

	for i, embed := range embeds {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d of %d", i+1, len(embeds)),
		}
	}

	return ctx.Kit.NewPaginate(ctx.Message.ChannelID, ctx.Message.Author.ID, embeds, time.Minute*10)
}

func (*ToneTags) Create(ctx *route.MessageContext) error {

	tag := ctx.Arguments["tag"].(string)
	tag = sanitiseShorthand(tag)

	desc := ctx.Arguments["description"].(string)

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

func (*ToneTags) Delete(ctx *route.MessageContext) error {
	tag := ctx.Arguments["tag"].(string)
	tag = sanitiseShorthand(tag)

	dbTag := &db.ToneTag{
		Shorthand: tag,
	}
	err := dbTag.Delete()
	if err != nil {
		return err
	}
	return ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "✅")
}
