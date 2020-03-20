package bete

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/yi-jiayu/ted"
)

const (
	callbackRefresh       = "refresh"
	callbackResend        = "resend"
	callbackAddFavourite  = "af"
	callbackSaveFavourite = "sf"
)

func (b Bete) HandleCallbackQuery(ctx context.Context, q *ted.CallbackQuery) {
	sentrySetUser(ctx, q.From.ID)

	var data CallbackData
	err := json.Unmarshal([]byte(q.Data), &data)
	if err != nil {
		return
	}
	switch data.Type {
	case callbackRefresh:
		b.updateETAs(ctx, q, data.StopID, data.Filter)
	case callbackResend:
		b.resendETAs(ctx, q, data.StopID, data.Filter)
	case callbackAddFavourite:
		b.askForFavouriteQuery(ctx, q)
	case callbackSaveFavourite:
		b.saveFavouriteCallback(ctx, q, data)
	}
}

func (b Bete) updateETAs(ctx context.Context, q *ted.CallbackQuery, stop string, filter []string) {
	text, err := b.etaMessageText(ctx, stop, filter)
	if err != nil {
		captureError(ctx, err)
		return
	}
	editMessageText := ted.EditMessageTextRequest{
		ChatID:      q.Message.Chat.ID,
		MessageID:   q.Message.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkupP(stop, filter),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
		Text:            "ETAs updated!",
	}
	_, err = b.Telegram.Do(editMessageText)
	if err != nil && !ted.IsMessageNotModified(err) {
		captureError(ctx, errors.WithStack(err))
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
}

func (b Bete) resendETAs(ctx context.Context, q *ted.CallbackQuery, stop string, filter []string) {
	text, err := b.etaMessageText(ctx, stop, filter)
	if err != nil {
		captureError(ctx, err)
		return
	}
	sendMessage := ted.SendMessageRequest{
		ChatID:      q.Message.Chat.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop, filter),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
		Text:            "ETAs sent!",
	}
	_, err = b.Telegram.Do(sendMessage)
	if err != nil && !ted.IsMessageNotModified(err) {
		captureError(ctx, errors.WithStack(err))
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
}

func (b Bete) askForFavouriteQuery(ctx context.Context, q *ted.CallbackQuery) {
	sendMessage := ted.SendMessageRequest{
		ChatID:      q.Message.Chat.ID,
		Text:        stringAddFavouritePromptForQuery,
		ReplyMarkup: ted.ForceReply{},
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	}
	var err error
	_, err = b.Telegram.Do(sendMessage)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
}

func (b Bete) saveFavouriteCallback(ctx context.Context, q *ted.CallbackQuery, data CallbackData) {
	query := data.Query
	if name := data.Name; name != "" {
		userID := q.From.ID
		err := b.Favourites.Put(userID, name, query.Canonical())
		if err != nil {
			captureError(ctx, err)
			return
		}
		favourites, err := b.Favourites.List(userID)
		if err != nil {
			captureError(ctx, err)
			return
		}
		showFavourites := ted.SendMessageRequest{
			ChatID:      q.Message.Chat.ID,
			Text:        fmt.Sprintf("Added the query %q to your favourites as %q!", query.Canonical(), data.Name),
			ReplyMarkup: showFavouritesReplyMarkup(favourites),
		}
		_, err = b.Telegram.Do(showFavourites)
		if err != nil {
			captureError(ctx, err)
		}
	} else {
		promptForName := ted.SendMessageRequest{
			ChatID:      q.Message.Chat.ID,
			Text:        fmt.Sprintf(stringAddFavouritePromptForName, query.Canonical()),
			ReplyMarkup: ted.ForceReply{},
		}
		_, err := b.Telegram.Do(promptForName)
		if err != nil {
			captureError(ctx, err)
		}
	}
	var err error
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	}
	removeButtons := ted.EditMessageReplyMarkupRequest{
		ChatID:    q.Message.Chat.ID,
		MessageID: q.Message.ID,
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		captureError(ctx, err)
	}
	_, err = b.Telegram.Do(removeButtons)
	if err != nil {
		captureError(ctx, err)
	}
}
