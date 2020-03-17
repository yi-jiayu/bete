package bete

import (
	"encoding/json"
	"log"

	"github.com/yi-jiayu/ted"
)

func (b Bete) HandleCallbackQuery(q *ted.CallbackQuery) {
	var data CallbackData
	err := json.Unmarshal([]byte(q.Data), &data)
	if err != nil {
		return
	}
	text, err := b.etaMessageText(data.StopID, data.Filter)
	if err != nil {
		log.Printf("error generating eta message text: %v", err)
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
		log.Printf("error making editMessageText request: %v", err)
	}
	_, err = b.Telegram.Do(answerCallbackQuery)
	if err != nil {
		log.Printf("error making answerCallbackQuery request: %v", err)
	}
}
