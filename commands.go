package bete

import (
	"context"
	"fmt"
	"regexp"

	"github.com/yi-jiayu/ted"
)

const (
	commandStart      = "start"
	commandFavourites = "favourites"
	commandAbout      = "about"
	commandVersion    = "version"
	commandETA        = "eta"
	commandTour       = "tour"
)

func (b Bete) HandleCommand(ctx context.Context, m *ted.Message, cmd, args string) {
	switch cmd {
	case commandStart:
		b.handleStartCommand(ctx, m)
	case commandFavourites:
		b.handleFavouritesCommand(ctx, m)
	case commandAbout:
		fallthrough
	case commandVersion:
		b.handleAboutCommand(ctx, m)
	case commandETA:
		b.handleETACommand(ctx, m, args)
	case commandTour:
		b.handleTourCommand(ctx, m)
	default:
		if isBusStopCodeCommand(cmd) {
			b.handleBusStopCodeCommand(ctx, m, cmd)
			cmd = "code"
		} else {
			captureMessage(ctx, "invalid command")
			b.handleInvalidCommand(ctx, m)
		}
		return
	}
	commandsTotal.WithLabelValues(cmd).Inc()
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

func isBusStopCodeCommand(command string) bool {
	match, err := regexp.MatchString(`\d{5}`, command)
	if err != nil {
		return false
	}
	return match
}

func (b Bete) handleBusStopCodeCommand(ctx context.Context, m *ted.Message, stopID string) {
	text, err := b.etaMessageText(ctx, stopID, nil, FormatSummary)
	if err != nil {
		captureError(ctx, err)
		return
	}
	req := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stopID, nil, FormatSummary),
	}
	b.send(ctx, req)
}

func (b Bete) handleAboutCommand(ctx context.Context, m *ted.Message) {
	req := ted.SendMessageRequest{
		ChatID:    m.Chat.ID,
		Text:      fmt.Sprintf(stringAboutMessage, b.Version, b.Version),
		ParseMode: "HTML",
	}
	b.send(ctx, req)
}

func (b Bete) handleStartCommand(ctx context.Context, m *ted.Message) {
	reply := ted.SendMessageRequest{
		ChatID: m.Chat.ID,
		Text:   fmt.Sprintf(stringWelcomeMessage, m.From.FirstName),
		ReplyMarkup: ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text: "Take the tour!",
						CallbackData: CallbackData{
							Type: callbackTour,
							Name: tourSectionStart,
						}.Encode(),
					},
				},
				{
					{
						Text: "About Bus Eta Bot",
						CallbackData: CallbackData{
							Type: callbackAbout,
						}.Encode(),
					},
				},
			},
		},
	}
	b.send(ctx, reply)
}

func (b Bete) handleTourCommand(ctx context.Context, m *ted.Message) {
	start := tour[tourSectionStart]
	reply := ted.SendMessageRequest{
		ChatID:      m.Chat.ID,
		Text:        start.Text,
		ReplyMarkup: tourReplyMarkup(start),
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

func (b Bete) handleInvalidCommand(ctx context.Context, m *ted.Message) {
	reply := ted.SendMessageRequest{
		ChatID: m.Chat.ID,
		Text:   stringInvalidCommand,
	}
	b.send(ctx, reply)
}
