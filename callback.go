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
	editMessageText := ted.EditMessageTextRequest{
		ChatID:      q.Message.Chat.ID,
		MessageID:   q.Message.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkupP(data.StopID, data.Filter),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
		Text:            "ETAs updated!",
	}
	_, err = b.Telegram.Do(editMessageText)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		captureError(ctx, errors.WithStack(err))
	}
}
