package bete

import (
	"strings"

	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleMessage(m *ted.Message) error {
	return b.HandleTextMessage(m)
}

func (b Bete) HandleTextMessage(m *ted.Message) error {
	parts := strings.Fields(m.Text)
	if len(parts) == 0 {
		return nil
	}
	stop, filter := parts[0], parts[1:]
	return b.SendETAMessage(m.Chat.ID, stop, filter)
}
