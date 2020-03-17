package bete

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleCallbackQuery(ctx context.Context, q *ted.CallbackQuery) {
	var data CallbackData
	err := json.Unmarshal([]byte(q.Data), &data)
	if err != nil {
		return
	}
	text, err := b.etaMessageText(ctx, data.StopID, data.Filter)
	if err != nil {
		captureError(ctx, err)
		return
	}
	switch data.Type {
	case "refresh":
		b.updateETAs(ctx, q, text, data.StopID, data.Filter)
	case "resend":
		b.resendETAs(ctx, q, text, data.StopID, data.Filter)
	}
}

func (b Bete) updateETAs(ctx context.Context, q *ted.CallbackQuery, text, stop string, filter []string) {
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
	_, err := b.Telegram.Do(editMessageText)
	if err != nil && !ted.IsMessageNotModified(err) {
		captureError(ctx, errors.WithStack(err))
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
}

func (b Bete) resendETAs(ctx context.Context, q *ted.CallbackQuery, text, stop string, filter []string) {
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
	_, err := b.Telegram.Do(sendMessage)
	if err != nil && !ted.IsMessageNotModified(err) {
		captureError(ctx, errors.WithStack(err))
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
}
