package bete

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleMessage(ctx context.Context, m *ted.Message) {
	sentrySetUser(ctx, m.From.ID)

	if m.Text == "" {
		// Ignore non-text messages.
		return
	}
	if cmd, args := m.CommandAndArgs(); cmd != "" {
		b.HandleCommand(ctx, m, cmd, args)
		return
	}
	if m.ReplyToMessage != nil {
		b.HandleReply(ctx, m)
		return
	}
	b.HandleTextMessage(ctx, m)
}

func (b Bete) HandleTextMessage(ctx context.Context, m *ted.Message) {
	var query string
	if favourite := b.Favourites.Find(m.From.ID, m.Text); favourite != "" {
		query = favourite
	} else {
		query = m.Text
		if len(query) > MaxQueryLength {
			_, err := b.Telegram.Do(ted.SendMessageRequest{
				ChatID: m.Chat.ID,
				Text:   ETAQueryTooLong,
			})
			if err != nil {
				captureError(ctx, err)
			}
			return
		}
	}
	if valid := validQueryRegexp.MatchString(query); !valid {
		return
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
	case "about":
		fallthrough
	case "version":
		b.handleAboutCommand(ctx, m)
	}
}

func (b Bete) handleAboutCommand(ctx context.Context, m *ted.Message) {
	req := ted.SendMessageRequest{
		ChatID:           m.Chat.ID,
		Text:             "Bus Eta Bot " + b.Version,
		ReplyToMessageID: m.ID,
	}
	_, err := b.Telegram.Do(req)
	if err != nil {
		captureError(ctx, err)
		return
	}
	return
}

func (b Bete) handleFavouritesCommand(ctx context.Context, m *ted.Message) {
	var req ted.Request
	if m.Chat.Type != "private" {
		req = ted.SendMessageRequest{
			ChatID:           m.Chat.ID,
			Text:             "Sorry, you can only manage your favourites in a private chat.",
			ReplyToMessageID: m.ID,
		}
	} else {
		req = ted.SendMessageRequest{
			ChatID:      m.Chat.ID,
			Text:        "What would you like to do?",
			ReplyMarkup: manageFavouritesReplyMarkup(),
		}
	}
	_, err := b.Telegram.Do(req)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
		return
	}
}

func getFavouriteQuery(text string) string {
	var query string
	n, err := fmt.Sscanf(text, AddFavouritePromptForName, &query)
	if err != nil {
		return ""
	}
	if n != 1 {
		return ""
	}
	return query
}

func (b Bete) HandleReply(ctx context.Context, m *ted.Message) {
	if m.ReplyToMessage.Text == AddFavouritePromptForQuery {
		b.addFavouriteSuggestName(ctx, m)
	} else if query := getFavouriteQuery(m.ReplyToMessage.Text); query != "" {
		b.addFavouriteFinish(ctx, m, query)
	}
}

func (b Bete) addFavouriteSuggestName(ctx context.Context, m *ted.Message) {
	query, err := ParseQuery(m.Text)
	if err != nil {
		reportError := ted.SendMessageRequest{
			ChatID:      m.Chat.ID,
			Text:        AddFavouriteReportQueryInvalid,
			ReplyMarkup: ted.ForceReply{},
		}
		askAgain := ted.SendMessageRequest{
			ChatID:      m.Chat.ID,
			Text:        AddFavouritePromptForQuery,
			ReplyMarkup: ted.ForceReply{},
		}
		var err error
		_, err = b.Telegram.Do(reportError)
		if err != nil {
			captureError(ctx, errors.WithStack(err))
			return
		}
		_, err = b.Telegram.Do(askAgain)
		if err != nil {
			captureError(ctx, errors.WithStack(err))
			return
		}
		return
	}
	var description string
	stop, err := b.BusStops.Find(query.Stop)
	switch {
	case err == ErrNotFound:
	case err != nil:
		captureError(ctx, err)
	default:
		description = stop.Description
	}
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        fmt.Sprintf(AddFavouriteSuggestName, query.Canonical()),
		ReplyMarkup: addFavouriteSuggestNameMarkup(query, description),
	}
	_, err = b.Telegram.Do(req)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
}

func (b Bete) addFavouriteFinish(ctx context.Context, m *ted.Message, query string) {
	name := m.Text
	userID := m.From.ID
	err := b.Favourites.Put(userID, name, query)
	if err != nil {
		captureError(ctx, err)
		return
	}
	favourites, err := b.Favourites.List(userID)
	if err != nil {
		captureError(ctx, err)
		return
	}
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        "New favourite added!",
		ReplyMarkup: showFavouritesReplyMarkup(favourites),
	}
	_, err = b.Telegram.Do(req)
	if err != nil {
		captureError(ctx, err)
		return
	}
}
