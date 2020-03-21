package bete

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yi-jiayu/ted"
)

func TestBete_HandleCallbackQuery_Refresh(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	editMessageText := ted.EditMessageTextRequest{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkupP(stop.ID, filter),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            "ETAs updated!",
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(editMessageText).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:   callbackRefresh,
				StopID: stop.ID,
				Filter: filter,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCallbackQuery_Resend(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	sendMessage := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            "ETAs sent!",
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(sendMessage).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:   callbackResend,
				StopID: stop.ID,
				Filter: filter,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCallbackQuery_AddFavourite(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	sendMessage := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        stringAddFavouritePromptForQuery,
		ReplyMarkup: ted.ForceReply{},
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(sendMessage).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackAddFavourite,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_saveFavouriteCallback_WithName(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	name := "Opp Tropicana Condo"
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	favourites := []string{"Home", name}
	showFavourites := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringAddFavouriteAdded, query.Canonical(), name),
		ReplyMarkup: showFavouritesReplyMarkup(favourites),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}
	removeButtons := ted.EditMessageReplyMarkupRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Put(userID, name, query.Canonical())
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(favourites, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(showFavourites).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(removeButtons).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:  callbackSaveFavourite,
				Query: &query,
				Name:  name,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_saveFavouriteCallback_WithoutName(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	promptForName := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringAddFavouritePromptForName, query.Canonical()),
		ReplyMarkup: ted.ForceReply{},
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}
	removeButtons := ted.EditMessageReplyMarkupRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(promptForName).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(removeButtons).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:  callbackSaveFavourite,
				Query: &query,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_deleteFavouritesCallback_ListError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            "Something went wrong!",
		CacheTime:       60,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(nil, Error("some error"))
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackDeleteFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_deleteFavouritesCallback_NoFavourites(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	editMessage := ted.EditMessageTextRequest{
		Text:        stringDeleteFavouritesNoFavourites,
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: deleteFavouritesReplyMarkupP(nil),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(nil, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(editMessage).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackDeleteFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_deleteFavouritesCallback(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	favourites := []string{"Home", "Work"}
	editMessage := ted.EditMessageTextRequest{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        stringDeleteFavouritesChoose,
		ReplyMarkup: deleteFavouritesReplyMarkupP(favourites),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(favourites, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(editMessage).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackDeleteFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_deleteFavouriteCallback(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	favouriteToDelete := "Work"
	remainingFavourites := []string{"Home"}
	editMessage := ted.EditMessageTextRequest{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        stringDeleteFavouritesChoose,
		ReplyMarkup: deleteFavouritesReplyMarkupP(remainingFavourites),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            stringDeleteFavouriteDeleted,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Delete(userID, favouriteToDelete).Return(nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(remainingFavourites, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(editMessage).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackDeleteFavourite,
				Name: favouriteToDelete,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_deleteFavouriteCallback_NoFavouritesRemaining(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	favouriteToDelete := "Work"
	editMessage := ted.EditMessageTextRequest{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        stringDeleteFavouritesNoFavouritesLeft,
		ReplyMarkup: deleteFavouritesReplyMarkupP(nil),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            stringDeleteFavouriteDeleted,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Delete(userID, favouriteToDelete).Return(nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(nil, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(editMessage).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackDeleteFavourite,
				Name: favouriteToDelete,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_deleteFavouriteCallback_DeleteError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	favouriteToDelete := "Work"
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            "Something went wrong!",
		CacheTime:       60,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Delete(userID, favouriteToDelete).Return(Error("some error"))
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackDeleteFavourite,
				Name: favouriteToDelete,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_deleteFavouriteCallback_ListError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	favouriteToDelete := "Work"
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            "Something went wrong!",
		CacheTime:       60,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Delete(userID, favouriteToDelete).Return(nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(nil, Error("some error"))
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackDeleteFavourite,
				Name: favouriteToDelete,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}
