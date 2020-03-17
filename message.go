package bete

import (
	"context"
	"log"
	"strings"

	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleMessage(ctx context.Context, m *ted.Message) {
	b.HandleTextMessage(ctx, m)
}

func (b Bete) HandleTextMessage(ctx context.Context, m *ted.Message) {
	parts := strings.Fields(m.Text)
	if len(parts) == 0 {
		return
	}
	stop, filter := parts[0], parts[1:]
	err := b.SendETAMessage(m.Chat.ID, stop, filter)
	if err != nil {
		log.Printf("error sending eta message: %v", err)
	}
}
