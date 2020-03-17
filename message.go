package bete

import (
	"context"
	"strings"

	"github.com/pkg/errors"
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
	text, err := b.etaMessageText(ctx, stop, filter)
	if err != nil {
		captureError(ctx, err)
		return
	}
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop, filter),
	}
	_, err = b.Telegram.Do(req)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
		return
	}
}
