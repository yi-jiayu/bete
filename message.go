package bete

import (
	"context"
	"fmt"

	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleMessage(ctx context.Context, m *ted.Message) {
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
	var query Query
	if favourite := b.Favourites.Find(m.From.ID, m.Text); favourite != "" {
		query, _ = ParseQuery(favourite)
	} else {
		var err error
		query, err = ParseQuery(m.Text)
		if err != nil {
			if err != ErrQueryDoesNotStartWithBusStopCode {
				b.reportInvalidQuery(ctx, m.Chat.ID, err)
			}
			return
		}
	}
	text, err := b.etaMessageText(ctx, query.Stop, query.Filter, FormatSummary)
	if err != nil {
		captureError(ctx, err)
		return
	}
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(query.Stop, query.Filter, FormatSummary),
	}
	b.send(ctx, req)
}

func (b Bete) reportInvalidQuery(ctx context.Context, chatID int64, err error) {
	var text string
	switch err {
	case ErrQueryDoesNotStartWithBusStopCode:
		text = stringQueryShouldStartWithBusStopCode
	case ErrQueryContainsInvalidCharacters:
		text = stringQueryContainsInvalidCharacters
	case ErrQueryTooLong:
		text = stringQueryTooLong
		// I want to know if anyone actually runs into this error.
		captureError(ctx, err)
	}
	b.send(ctx, ted.SendMessageRequest{
		ChatID: chatID,
		Text:   text,
	})
}

func (b Bete) HandleCommand(ctx context.Context, m *ted.Message, cmd, args string) {
	switch cmd {
	case "start":
		b.handleStartCommand(ctx, m)
	case "favourites":
		b.handleFavouritesCommand(ctx, m)
	case "about":
		fallthrough
	case "version":
		b.handleAboutCommand(ctx, m)
	case "eta":
		b.handleETACommand(ctx, m, args)
	case callbackTour:
		b.handleTourCommand(ctx, m)
	}
}

func (b Bete) handleETACommand(ctx context.Context, m *ted.Message, args string) {
	if args == "" {
		b.handleETACommandWithoutArgs(ctx, m)
		return
	}
	query, err := ParseQuery(args)
	if err != nil {
		b.reportInvalidQuery(ctx, m.Chat.ID, err)
		return
	}
	text, err := b.etaMessageText(ctx, query.Stop, query.Filter, FormatSummary)
	if err != nil {
		captureError(ctx, err)
		return
	}
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(query.Stop, query.Filter, FormatSummary),
	}
	b.send(ctx, req)
}

func (b Bete) handleETACommandWithoutArgs(ctx context.Context, m *ted.Message) {
	reply := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        stringETACommandPrompt,
		ParseMode:   "HTML",
		ReplyMarkup: ted.ForceReply{},
	}
	b.send(ctx, reply)
}

func (b Bete) handleAboutCommand(ctx context.Context, m *ted.Message) {
	req := ted.SendMessageRequest{
		ChatID:           m.Chat.ID,
		Text:             "Bus Eta Bot " + b.Version,
		ReplyToMessageID: m.ID,
	}
	b.send(ctx, req)
}

func (b Bete) handleStartCommand(ctx context.Context, m *ted.Message) {
	reply := ted.SendMessageRequest{
		ChatID: m.Chat.ID,
		Text:   fmt.Sprintf(stringWelcomeMessage, m.From.FirstName),
	}
	b.send(ctx, reply)
}

func (b Bete) handleTourCommand(ctx context.Context, m *ted.Message) {
	reply := ted.SendMessageRequest{
		ChatID: m.Chat.ID,
		Text:   stringTourStart,
		ReplyMarkup: ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text: "Next: " + stringTourTitleETAQueries,
						CallbackData: CallbackData{
							Type: callbackTour,
							Name: "eta_queries",
						}.Encode(),
					},
				},
			},
		},
	}
	b.send(ctx, reply)
}

func (b Bete) handleFavouritesCommand(ctx context.Context, m *ted.Message) {
	var req ted.Request
	if m.Chat.Type != "private" {
		req = ted.SendMessageRequest{
			ChatID:           m.Chat.ID,
			Text:             stringFavouritesOnlyPrivateChat,
			ReplyToMessageID: m.ID,
		}
	} else {
		req = ted.SendMessageRequest{
			ChatID:      m.Chat.ID,
			Text:        stringFavouritesChooseAction,
			ReplyMarkup: favouritesReplyMarkup(),
		}
	}
	b.send(ctx, req)
}

func getFavouriteQuery(text string) string {
	var query string
	n, err := fmt.Sscanf(text, stringAddFavouritePromptForName, &query)
	if err != nil {
		return ""
	}
	if n != 1 {
		return ""
	}
	return query
}

// HandleReply handles messages which are replies.
//
// When matching against the reply text, note that formatting markup in the original message will not be present.
func (b Bete) HandleReply(ctx context.Context, m *ted.Message) {
	if m.ReplyToMessage.Text == stringAddFavouritePromptForQuery {
		b.addFavouriteSuggestName(ctx, m)
	} else if query := getFavouriteQuery(m.ReplyToMessage.Text); query != "" {
		b.addFavouriteFinish(ctx, m, query)
	} else if m.ReplyToMessage.Text == stringETACommandPrompt {
		b.handleETACommand(ctx, m, m.Text)
	} else {
		b.send(ctx, ted.SendMessageRequest{
			ChatID: m.Chat.ID,
			Text:   "Sorry, I forgot what we were talking about.",
		})
	}
}

func (b Bete) addFavouriteSuggestName(ctx context.Context, m *ted.Message) {
	query, err := ParseQuery(m.Text)
	if err != nil {
		b.reportInvalidQuery(ctx, m.Chat.ID, err)
		askAgain := ted.SendMessageRequest{
			ChatID:      m.Chat.ID,
			Text:        stringAddFavouritePromptForQuery,
			ReplyMarkup: ted.ForceReply{},
		}
		b.send(ctx, askAgain)
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
		Text:        fmt.Sprintf(stringAddFavouriteSuggestName, query.Canonical()),
		ReplyMarkup: addFavouriteSuggestNameMarkup(query, description),
	}
	b.send(ctx, req)
}

func (b Bete) addFavouriteFinish(ctx context.Context, m *ted.Message, query string) {
	name := m.Text
	userID := m.From.ID
	err := b.Favourites.Put(userID, name, query)
	if err != nil {
		captureError(ctx, err)
		b.send(ctx, ted.SendMessageRequest{
			ChatID: m.Chat.ID,
			Text:   stringErrorSorry,
		})
		return
	}
	favourites, err := b.Favourites.List(userID)
	if err != nil {
		captureError(ctx, err)
		b.send(ctx, ted.SendMessageRequest{
			ChatID: m.Chat.ID,
			Text:   stringErrorSorry,
		})
		return
	}
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        fmt.Sprintf(stringAddFavouriteAdded, query, name),
		ReplyMarkup: showFavouritesReplyMarkup(favourites),
	}
	b.send(ctx, req)
}
