package bete

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleCallbackQuery(ctx context.Context, q *ted.CallbackQuery) {
	sentrySetUser(ctx, q.From.ID)

	var data CallbackData
	err := json.Unmarshal([]byte(q.Data), &data)
	if err != nil {
		return
	}
	switch data.Type {
	case "refresh":
		b.updateETAs(ctx, q, data.StopID, data.Filter)
	case "resend":
		b.resendETAs(ctx, q, data.StopID, data.Filter)
	case "add_favourite":
		b.askForFavouriteQuery(ctx, q)
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
		Text:        AddFavouritePromptForQuery,
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
