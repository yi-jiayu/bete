package bete

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yi-jiayu/ted"
)

// Callback query types
const (
	callbackRefresh          = "refresh"
	callbackResend           = "resend"
	callbackAddFavourite     = "af"
	callbackSaveFavourite    = "sf"
	callbackFavourites       = "favourites"
	callbackDeleteFavourites = "delete_favourites"
	callbackDeleteFavourite  = "delete_favourite"
	callbackShowFavourites   = "show_favourites"
	callbackHideFavourites   = "hide_favourites"
	callbackTour             = "tour"
	callbackAbout            = "about"
	callbackNearbyETA        = "nearby_eta"
)

func (b Bete) HandleCallbackQuery(ctx context.Context, q *ted.CallbackQuery) {
	var data CallbackData
	err := json.Unmarshal([]byte(q.Data), &data)
	if err != nil {
		captureMessage(ctx, "unrecogised callback")
		return
	}
	switch data.Type {
	case callbackRefresh:
		b.updateETAs(ctx, q, data)
	case callbackResend:
		b.resendETAs(ctx, q, data)
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
	case callbackHideFavourites:
		b.hideFavouritesCallback(ctx, q)
	case callbackFavourites:
		b.favouritesCallback(ctx, q)
	case callbackTour:
		b.tourCallback(ctx, q, data)
	case callbackAbout:
		b.aboutCallback(ctx, q)
	case callbackNearbyETA:
		b.resendETAs(ctx, q, data)
	default:
		captureMessage(ctx, "unrecogised callback")
		return
	}
	callbackQueriesTotal.WithLabelValues(data.Type).Inc()
}

func (b Bete) answerCallbackQueryError(ctx context.Context, q *ted.CallbackQuery, err error) {
	captureError(ctx, err)
	b.send(ctx, ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
		Text:            stringSomethingWentWrong,
		CacheTime:       60,
	})
}

func (b Bete) updateETAs(ctx context.Context, q *ted.CallbackQuery, data CallbackData) {
	text, err := b.etaMessageText(ctx, data.StopID, data.Filter, data.Format)
	if err != nil {
		captureError(ctx, err)
		return
	}
	editMessageText := ted.EditMessageTextRequest{
		Text:      text,
		ParseMode: "HTML",
	}
	if q.Message != nil {
		editMessageText.ChatID = q.Message.Chat.ID
		editMessageText.MessageID = q.Message.ID
		editMessageText.ReplyMarkup = etaMessageReplyMarkupP(data.StopID, data.Filter, data.Format)
	} else {
		editMessageText.InlineMessageID = q.InlineMessageID
		editMessageText.ReplyMarkup = inlineETAMessageReplyMarkupP(data.StopID, data.Format)
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
		Text:            stringRefreshETAsUpdated,
	}
	b.send(ctx, editMessageText)
	b.send(ctx, answerCallbackQuery)
}

func (b Bete) resendETAs(ctx context.Context, q *ted.CallbackQuery, data CallbackData) {
	text, err := b.etaMessageText(ctx, data.StopID, data.Filter, data.Format)
	if err != nil {
		captureError(ctx, err)
		return
	}
	sendMessage := ted.SendMessageRequest{
		ChatID:      q.Message.Chat.ID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(data.StopID, data.Filter, data.Format),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
		Text:            stringResendETAsSent,
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
			b.answerCallbackQueryError(ctx, q, err)
			return
		}
		favourites, err := b.Favourites.List(userID)
		if err != nil {
			b.answerCallbackQueryError(ctx, q, err)
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
			Text:        stringShowFavouritesShowing,
			ReplyMarkup: showFavouritesReplyMarkup(favourites),
		}
	}
	b.send(ctx, req)
	b.send(ctx, ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	})
}

func (b Bete) hideFavouritesCallback(ctx context.Context, q *ted.CallbackQuery) {
	hideKeyboard := ted.SendMessageRequest{
		ChatID:      q.Message.Chat.ID,
		Text:        stringHideFavouritesHiding,
		ReplyMarkup: ted.ReplyKeyboardRemove{},
	}
	b.send(ctx, hideKeyboard)
}

func (b Bete) favouritesCallback(ctx context.Context, q *ted.CallbackQuery) {
	editMessage := ted.EditMessageTextRequest{
		ChatID:      q.Message.Chat.ID,
		MessageID:   q.Message.ID,
		Text:        stringFavouritesChooseAction,
		ReplyMarkup: favouritesReplyMarkupP(),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	}
	b.send(ctx, editMessage)
	b.send(ctx, answerCallback)
}

func (b Bete) tourCallback(ctx context.Context, q *ted.CallbackQuery, data CallbackData) {
	section := tour[data.Name]
	reply := ted.SendMessageRequest{
		ChatID:      q.Message.Chat.ID,
		Text:        section.Text,
		ParseMode:   "HTML",
		ReplyMarkup: tourReplyMarkup(section),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	}
	b.send(ctx, reply)
	b.send(ctx, answerCallback)
}

func (b Bete) aboutCallback(ctx context.Context, q *ted.CallbackQuery) {
	aboutMessage := ted.SendMessageRequest{
		ChatID:    q.Message.Chat.ID,
		Text:      fmt.Sprintf(stringAboutMessage, b.Version, b.Version),
		ParseMode: "HTML",
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: q.ID,
	}
	b.send(ctx, aboutMessage)
	b.send(ctx, answerCallback)
}
