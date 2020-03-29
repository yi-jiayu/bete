package bete

import (
	"context"
	"fmt"
	"testing"

	"github.com/yi-jiayu/ted"
)

func TestBete_updateETAs(t *testing.T) {

	tests := []struct {
		name   string
		format Format
	}{
		{
			name:   "arriving bus summary",
			format: FormatSummary,
		},
		{
			name:   "arriving bus details",
			format: FormatDetails,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, finish := newMockBete(t)
			defer finish()

			stop := buildBusStop()
			filter := []string{"5", "24"}
			arrivals := buildDataMallBusArrival()
			chatID := randomInt64ID()
			messageID := randomID()
			callbackQueryID := randomStringID()
			text := must(FormatArrivals(ArrivalInfo{
				Stop:     stop,
				Time:     refTime,
				Services: arrivals.Services,
				Filter:   filter,
			}, tt.format)).(string)
			editMessageText := ted.EditMessageTextRequest{
				ChatID:      chatID,
				MessageID:   messageID,
				Text:        text,
				ParseMode:   "HTML",
				ReplyMarkup: etaMessageReplyMarkupP(stop.ID, filter, tt.format),
			}
			answerCallbackQuery := ted.AnswerCallbackQueryRequest{
				CallbackQueryID: callbackQueryID,
				Text:            stringRefreshETAsUpdated,
			}

			b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
			b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
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
						Format: tt.format,
					}.Encode(),
				},
			}
			b.HandleUpdate(context.Background(), update)
		})
	}
}

func TestBete_updateETAs_Inline(t *testing.T) {
	tests := []struct {
		name   string
		format Format
	}{
		{
			name:   "arriving bus summary",
			format: FormatSummary,
		},
		{
			name:   "arriving bus details",
			format: FormatDetails,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, finish := newMockBete(t)
			defer finish()

			stop := buildBusStop()
			filter := []string{"5", "24"}
			arrivals := buildDataMallBusArrival()
			inlineMessageID := randomStringID()
			callbackQueryID := randomStringID()
			text := must(FormatArrivals(ArrivalInfo{
				Stop:     stop,
				Time:     refTime,
				Services: arrivals.Services,
				Filter:   filter,
			}, tt.format)).(string)
			editMessageText := ted.EditMessageTextRequest{
				InlineMessageID: inlineMessageID,
				Text:            text,
				ParseMode:       "HTML",
				ReplyMarkup:     inlineETAMessageReplyMarkupP(stop.ID, tt.format),
			}
			answerCallbackQuery := ted.AnswerCallbackQueryRequest{
				CallbackQueryID: callbackQueryID,
				Text:            stringRefreshETAsUpdated,
			}

			b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
			b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
			b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
			b.Telegram.(*MockTelegram).EXPECT().Do(editMessageText).Return(ted.Response{}, nil)
			b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

			update := ted.Update{
				CallbackQuery: &ted.CallbackQuery{
					ID:              callbackQueryID,
					InlineMessageID: inlineMessageID,
					Data: CallbackData{
						Type:   callbackRefresh,
						StopID: stop.ID,
						Filter: filter,
						Format: tt.format,
					}.Encode(),
				},
			}
			b.HandleUpdate(context.Background(), update)
		})
	}
}

func TestBete_resendETAs(t *testing.T) {
	tests := []struct {
		name   string
		format Format
	}{
		{
			name:   "arriving bus summary",
			format: FormatSummary,
		},
		{
			name:   "arriving bus details",
			format: FormatDetails,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, finish := newMockBete(t)
			defer finish()

			stop := buildBusStop()
			filter := []string{"5", "24"}
			arrivals := buildDataMallBusArrival()
			chatID := randomInt64ID()
			messageID := randomID()
			callbackQueryID := randomStringID()
			text := must(FormatArrivals(ArrivalInfo{
				Stop:     stop,
				Time:     refTime,
				Services: arrivals.Services,
				Filter:   filter,
			}, tt.format)).(string)
			sendMessage := ted.SendMessageRequest{
				ChatID:      chatID,
				Text:        text,
				ParseMode:   "HTML",
				ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter, tt.format),
			}
			answerCallbackQuery := ted.AnswerCallbackQueryRequest{
				CallbackQueryID: callbackQueryID,
				Text:            stringResendETAsSent,
			}

			b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
			b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
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
						Format: tt.format,
					}.Encode(),
				},
			}
			b.HandleUpdate(context.Background(), update)
		})
	}
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

func TestBete_saveFavouriteCallback_WithName_PutError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	name := "Opp Tropicana Condo"
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            stringSomethingWentWrong,
		CacheTime:       60,
	}
	b.Favourites.(*MockFavouriteRepository).EXPECT().Put(userID, name, query.Canonical()).Return(Error("some error"))
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

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

func TestBete_saveFavouriteCallback_WithName_ListError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	name := "Opp Tropicana Condo"
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            stringSomethingWentWrong,
		CacheTime:       60,
	}
	b.Favourites.(*MockFavouriteRepository).EXPECT().Put(userID, name, query.Canonical())
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(nil, Error("some error"))
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

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
		Text:            stringSomethingWentWrong,
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
	removeButtons := ted.EditMessageReplyMarkupRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}
	showFavourites := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringDeleteFavouriteDeleted, favouriteToDelete),
		ReplyMarkup: showFavouritesReplyMarkup(remainingFavourites),
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Delete(userID, favouriteToDelete).Return(nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(remainingFavourites, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(removeButtons).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(showFavourites).Return(ted.Response{}, nil)

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
		Text:            stringSomethingWentWrong,
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
		Text:            stringSomethingWentWrong,
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

func TestBete_showFavouritesCallback_ListError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            stringSomethingWentWrong,
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
				Type: callbackShowFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_showFavouritesCallback_NoFavourites(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	editMessage := ted.EditMessageTextRequest{
		Text:      stringShowFavouritesNoFavourites,
		ChatID:    chatID,
		MessageID: messageID,
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
				Type: callbackShowFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_showFavouritesCallback(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	favourites := []string{"Home", "Work"}
	showFavourites := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        stringShowFavouritesShowing,
		ReplyMarkup: showFavouritesReplyMarkup(favourites),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(favourites, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(showFavourites).Return(ted.Response{}, nil)
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
				Type: callbackShowFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_hideFavouritesCallback(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	hideKeyboard := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        stringHideFavouritesHiding,
		ReplyMarkup: ted.ReplyKeyboardRemove{},
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(hideKeyboard).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackHideFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_favouritesCallback(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	callbackQueryID := randomStringID()
	messageID := randomID()
	chatID := randomInt64ID()
	editMessage := ted.EditMessageTextRequest{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        stringFavouritesChooseAction,
		ReplyMarkup: favouritesReplyMarkupP(),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(editMessage).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				ID:   messageID,
				From: &ted.User{ID: userID},
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackFavourites,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_tourCallback(t *testing.T) {
	section := tourSectionStart
	callbackQueryID := randomStringID()
	chatID := randomInt64ID()

	b, finish := newMockBete(t)
	defer finish()

	sendMessage := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        tour[section].Text,
		ParseMode:   "HTML",
		ReplyMarkup: tourReplyMarkup(tour[section]),
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(sendMessage).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackTour,
				Name: section,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_aboutCallback(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	version := randomStringID()
	callbackQueryID := randomStringID()
	chatID := randomInt64ID()
	about := ted.SendMessageRequest{
		ChatID:    chatID,
		Text:      fmt.Sprintf(stringAboutMessage, version, version),
		ParseMode: "HTML",
	}
	answerCallback := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

	b.Version = version
	b.Telegram.(*MockTelegram).EXPECT().Do(about).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallback).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type: callbackAbout,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}
