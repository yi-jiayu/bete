package bete

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yi-jiayu/ted"
)

const (
	callbackRefresh          = "refresh"
	callbackResend           = "resend"
	callbackAddFavourite     = "af"
	callbackSaveFavourite    = "sf"
	callbackDeleteFavourites = "delete_favourites"
	callbackDeleteFavourite  = "delete_favourite"
	callbackShowFavourites   = "show_favourites"
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
	case callbackDeleteFavourites:
		b.deleteFavouritesCallback(ctx, q)
	case callbackDeleteFavourite:
		b.deleteFavouriteCallback(ctx, q, data)
	case callbackShowFavourites:
		b.showFavouritesCallback(ctx, q)
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
	b.send(ctx, editMessageText)
	b.send(ctx, answerCallbackQuery)
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
	b.send(ctx, sendMessage)
	b.send(ctx, answerCallbackQuery)
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
	b.send(ctx, sendMessage)
	b.send(ctx, answerCallbackQuery)
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
			Text:        fmt.Sprintf(stringAddFavouriteAdded, query.Canonical(), data.Name),
			ReplyMarkup: showFavouritesReplyMarkup(favourites),
		}
		b.send(ctx, showFavourites)
	} else {
		promptForName := ted.SendMessageRequest{
			ChatID:      q.Message.Chat.ID,
			Text:        fmt.Sprintf(stringAddFavouritePromptForName, query.Canonical()),
			ReplyMarkup: ted.ForceReply{},
		}
		b.send(ctx, promptForName)
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	}
	removeButtons := ted.EditMessageReplyMarkupRequest{
		ChatID:    q.Message.Chat.ID,
		MessageID: q.Message.ID,
	}
	b.send(ctx, answerCallbackQuery)
	b.send(ctx, removeButtons)
}

func (b Bete) deleteFavouritesCallback(ctx context.Context, q *ted.CallbackQuery) {
	favourites, err := b.Favourites.List(q.From.ID)
	if err != nil {
		b.answerCallbackQueryError(ctx, q, err)
		return
	}
	var text string
	if len(favourites) == 0 {
		text = stringDeleteFavouritesNoFavourites
	} else {
		text = stringDeleteFavouritesChoose
	}
	b.send(ctx, ted.EditMessageTextRequest{
		Text:        text,
		ChatID:      q.Message.Chat.ID,
		MessageID:   q.Message.ID,
		ReplyMarkup: deleteFavouritesReplyMarkupP(favourites),
	})
	b.send(ctx, ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	})
}

func (b Bete) deleteFavouriteCallback(ctx context.Context, q *ted.CallbackQuery, data CallbackData) {
	userID := q.From.ID
	favouriteToDelete := data.Name
	err := b.Favourites.Delete(userID, favouriteToDelete)
	if err != nil {
		b.answerCallbackQueryError(ctx, q, err)
		return
	}
	remainingFavourites, err := b.Favourites.List(userID)
	if err != nil {
		b.answerCallbackQueryError(ctx, q, err)
		return
	}
	b.send(ctx, ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	})
	chatID := q.Message.Chat.ID
	b.send(ctx, ted.EditMessageReplyMarkupRequest{
		ChatID:    chatID,
		MessageID: q.Message.ID,
	})
	b.send(ctx, ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringDeleteFavouriteDeleted, favouriteToDelete),
		ReplyMarkup: showFavouritesReplyMarkup(remainingFavourites),
	})
}

func (b Bete) showFavouritesCallback(ctx context.Context, q *ted.CallbackQuery) {
	favourites, err := b.Favourites.List(q.From.ID)
	if err != nil {
		b.answerCallbackQueryError(ctx, q, err)
		return
	}
	var req ted.Request
	if len(favourites) == 0 {
		req = ted.EditMessageTextRequest{
			Text:      stringShowFavouritesNoFavourites,
			ChatID:    q.Message.Chat.ID,
			MessageID: q.Message.ID,
			ReplyMarkup: &ted.InlineKeyboardMarkup{
				InlineKeyboard: [][]ted.InlineKeyboardButton{
					{
						{
							Text:         stringFavouritesAddNew,
							CallbackData: CallbackData{Type: callbackAddFavourite}.Encode(),
						},
					},
				},
			},
		}
	} else {
		req = ted.SendMessageRequest{
			ChatID:      q.Message.Chat.ID,
			Text:        "Showing favourites keyboard",
			ReplyMarkup: showFavouritesReplyMarkup(favourites),
		}
	}
	b.send(ctx, req)
	b.send(ctx, ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	})
}

func (b Bete) answerCallbackQueryError(ctx context.Context, q *ted.CallbackQuery, err error) {
	captureError(ctx, err)
	b.send(ctx, ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
		Text:            "Something went wrong!",
		CacheTime:       60,
	})
}
