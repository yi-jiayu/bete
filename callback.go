package bete

import (
	"encoding/json"
	"log"

	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleCallbackQuery(q *ted.CallbackQuery) error {
	var data CallbackData
	err := json.Unmarshal([]byte(q.Data), &data)
	if err != nil {
		return err
	}
	text, err := b.etaMessageText(data.StopID, data.Filter)
	if err != nil {
		return err
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
		log.Printf("error making editMessageText request: %v", err)
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		log.Printf("error making answerCallbackQuery request: %v", err)
	}
	return nil
}
