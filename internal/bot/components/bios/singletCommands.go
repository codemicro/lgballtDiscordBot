package bios

import (
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
)

func (*Bios) SingletHelp(ctx *route.MessageContext) error {
	_, err := ctx.SendMessageEmbed(ctx.Message.ChannelID, helpEmbed)
	return err
}

func (b *Bios) SingletSetField(ctx *route.MessageContext) error {

	newValue := ctx.Arguments["newValue"].(string)
	field := ctx.Arguments["field"].(string)

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID

	return b.setBioField(bdt, field, newValue, false, ctx)
}

func (b *Bios) SingletClearField(ctx *route.MessageContext) error {

	field := ctx.Arguments["field"].(string)

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID

	return b.clearBioField(bdt, field, ctx)
}

func (b *Bios) SingletSetImage(ctx *route.MessageContext) error {

	imageURL := ctx.Arguments["imageURL"].(string)

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID

	return b.setBioImage(bdt, imageURL, false, ctx)
}

func (b *Bios) SingletClearImage(ctx *route.MessageContext) error {

	bdt := new(db.UserBio)
	bdt.UserId = ctx.Message.Author.ID

	return b.clearBioImage(bdt, ctx)
}