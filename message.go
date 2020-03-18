package bete

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleMessage(ctx context.Context, m *ted.Message) {
	if cmd, args := m.CommandAndArgs(); cmd != "" {
		b.HandleCommand(ctx, m, cmd, args)
		return
	}
	b.HandleTextMessage(ctx, m)
}

func (b Bete) HandleTextMessage(ctx context.Context, m *ted.Message) {
	var query string
	if favourite := b.Favourites.FindByUserAndText(m.From.ID, m.Text); favourite != "" {
		query = favourite
	} else {
		query = m.Text
	}
	parts := strings.Fields(query)
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

func (b Bete) HandleCommand(ctx context.Context, m *ted.Message, cmd, args string) {
	switch cmd {
	case "favourites":
		b.handleFavouritesCommand(ctx, m)
	}
}

func (b Bete) handleFavouritesCommand(ctx context.Context, m *ted.Message) {
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        "What would you like to do?",
		ReplyMarkup: favouritesReplyMarkup(),
	}
	_, err := b.Telegram.Do(req)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
		return
	}
}
