package bios

import (
	"fmt"
	"github.com/skwair/harmony"
)

func (b *Bios) HelpSystem(_ []string, m *harmony.Message) error {
	_, err := b.b.SendMessage(m.ChannelID, "TODO") // TODO
	return err
}

func (b *Bios) SetFieldSystem(command []string, m *harmony.Message) error {
	fmt.Println("SetFieldSystem")
	return nil
}

func (b *Bios) ClearFieldSystem(command []string, m *harmony.Message) error {
	fmt.Println("ClearFieldSystem")
	return nil
}

