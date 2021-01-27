package bios

import (
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/skwair/harmony"
	"strings"
)

func (b *Bios) HelpSystem(_ []string, m *harmony.Message) error {
	_, err := b.b.SendMessage(m.ChannelID, "TODO") // TODO
	return err
}

func (b *Bios) SetFieldSystem(command []string, m *harmony.Message) error {
	// Syntax: <member ID> <field name> <value>

	newValue := strings.Join(command[2:], " ")
	bdt := new(db.UserBio)
	bdt.UserId = m.Author.ID
	bdt.SysMemberID = command[0]

	return b.setBioField(bdt, command[1], newValue, true, m)
}

func (b *Bios) ClearFieldSystem(command []string, m *harmony.Message) error {
	// Syntax: <member ID> <field name>

	bdt := new(db.UserBio)
	bdt.UserId = m.Author.ID
	bdt.SysMemberID = command[0]

	return b.clearBioField(bdt, command[1], m)
}
