package bete

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/ted"
)

func TestBete_HandleTextMessage(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	userID := randomID()
	chatID := randomInt64ID()
	text := must(formatArrivalsSummary(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter, FormatSummary),
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().Find(gomock.Any(), gomock.Any()).Return("")
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: "96049 5 24",
		},
	}
	b.HandleUpdate(context.Background(), update)
}

// Do not respond when an incoming text message does not start with a 5-digit bus stop code (after favourites check).
func TestBete_HandleTextMessage_InvalidQuery(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	text := "Hello, World"
	chatID := randomInt64ID()

	b.Favourites.(*MockFavouriteRepository).EXPECT().Find(gomock.Any(), gomock.Any()).Return("")

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: randomID()},
			Chat: ted.Chat{ID: chatID},
			Text: text,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

// Display an error for queries longer than 32 bytes (an arbitrary limit, but chosen to stay within the 64-character limit for callback query data).
func TestBete_HandleTextMessage_LongQuery(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	query := "81111 155 134 135 135A 137 154 155 24 28 43 43e 70 70A 70M 76"
	chatID := randomInt64ID()
	reply := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   stringQueryTooLong,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Find(gomock.Any(), gomock.Any()).Return("")
	b.Telegram.(*MockTelegram).EXPECT().Do(reply)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: randomID()},
			Chat: ted.Chat{ID: chatID},
			Text: query,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleTextMessage_Favourite(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	userID := randomID()
	chatID := randomInt64ID()
	messageText := "SUTD"
	replyText := must(formatArrivalsSummary(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        replyText,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter, FormatSummary),
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().Find(userID, messageText).Return("96049 5 24")
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: messageText,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_addFavouriteSuggestName(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	stop := buildBusStop()
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringAddFavouriteSuggestName, query.Canonical()),
		ReplyMarkup: addFavouriteSuggestNameMarkup(query, stop.Description),
	}

	b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: stringAddFavouritePromptForQuery,
			},
			Text: query.Canonical(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_addFavouriteSuggestName_BusStopNotFound(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	stop := buildBusStop()
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringAddFavouriteSuggestName, query.Canonical()),
		ReplyMarkup: addFavouriteSuggestNameMarkup(query, ""),
	}

	b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(BusStop{}, ErrNotFound)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: stringAddFavouritePromptForQuery,
			},
			Text: query.Canonical(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_addFavouriteFinish(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	query := "96049 5 24"
	favourites := []string{"Home", "SUTD"}
	name := "SUTD"
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringAddFavouriteAdded, query, name),
		ReplyMarkup: showFavouritesReplyMarkup(favourites),
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Put(userID, name, query).Return(nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(favourites, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: fmt.Sprintf(stringAddFavouritePromptForName, query),
			},
			Text: name,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_addFavouriteFinish_PutError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	query := "96049 5 24"
	name := "SUTD"
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   stringErrorSorry,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Put(userID, name, query).Return(Error("some error"))
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: fmt.Sprintf(stringAddFavouritePromptForName, query),
			},
			Text: name,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_addFavouriteFinish_ListError(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	query := "96049 5 24"
	name := "SUTD"
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   stringErrorSorry,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Put(userID, name, query).Return(nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(nil, Error("some error"))
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: fmt.Sprintf(stringAddFavouritePromptForName, query),
			},
			Text: name,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleReply_AddFavourite_HandleInvalidQuery(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageText := `Invalid Query: !@#$%^&*"`
	reportError := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   stringQueryShouldStartWithBusStopCode,
	}
	askAgain := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        stringAddFavouritePromptForQuery,
		ReplyMarkup: ted.ForceReply{},
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(reportError).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(askAgain).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: stringAddFavouritePromptForQuery,
			},
			Text: messageText,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleReply_etaCommandArgs(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	query := Query{Stop: stop.ID, Filter: []string{"5", "24"}}
	arrivals := buildDataMallBusArrival()
	userID := randomID()
	messageID := randomID()
	chatID := randomInt64ID()
	text := must(formatArrivalsSummary(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   query.Filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, query.Filter, FormatSummary),
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			ID:   messageID,
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: query.Canonical(),
			ReplyToMessage: &ted.Message{
				Text: stringETACommandPrompt,
			},
		},
	}
	b.HandleUpdate(context.Background(), update)

}

func TestBete_HandleReply_etaCommandArgsInvalid(t *testing.T) {
	tests := []struct {
		name         string
		args         string
		expectedText string
	}{
		{
			name:         "does not start with a bus stop code",
			args:         "ABCDE 5 24",
			expectedText: stringQueryShouldStartWithBusStopCode,
		},
		{
			name:         "contains invalid characters",
			args:         "12345 5! '2'",
			expectedText: stringQueryContainsInvalidCharacters,
		},
		{
			name:         "too long",
			args:         "12345 24 28 43 70 76 134 135 137 154 155",
			expectedText: stringQueryTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, finish := newMockBete(t)
			defer finish()

			userID := randomID()
			messageID := randomID()
			chatID := randomInt64ID()
			req := ted.SendMessageRequest{
				ChatID: chatID,
				Text:   tt.expectedText,
			}

			b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

			update := ted.Update{
				Message: &ted.Message{
					ID:   messageID,
					From: &ted.User{ID: userID},
					Chat: ted.Chat{ID: chatID},
					Text: tt.args,
					ReplyToMessage: &ted.Message{
						Text: stringETACommandPrompt,
					},
				},
			}
			b.HandleUpdate(context.Background(), update)
		})
	}
}

func TestBete_HandleReply_noMatch(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   "Sorry, I forgot what we were talking about.",
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: "anything",
			ReplyToMessage: &ted.Message{
				Text: "forgotten message",
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func Test_getFavouriteQuery(t *testing.T) {
	text := fmt.Sprintf(stringAddFavouritePromptForName, "96049 5 24")
	assert.Equal(t, "96049 5 24", getFavouriteQuery(text))
	assert.Equal(t, "", getFavouriteQuery("invalid"))
}

func TestBete_HandleReply_locationQuery(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	stops := []NearbyBusStop{
		{
			BusStop: BusStop{
				ID:          "01339",
				Description: "Bef Crawford Bridge",
				RoadName:    "Crawford St",
				Location:    Location{Latitude: 1.307746, Longitude: 103.864263},
			},
			Distance: 0.11356564947243729,
		},
		{
			BusStop: BusStop{
				ID:          "07371",
				Description: "Aft Kallang Rd",
				RoadName:    "Lavender St",
				Location:    Location{Latitude: 1.309508, Longitude: 103.863501},
			},
			Distance: 0.21676780485189698,
		},
	}
	var lat, lon float32 = 1.307574, 103.863256
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   stringLocationNearby,
	}
	venues := make([]ted.Request, len(stops))
	for i := range stops {
		venues[i] = ted.SendVenueRequest{
			ChatID:    chatID,
			Latitude:  stops[i].BusStop.Location.Latitude,
			Longitude: stops[i].BusStop.Location.Longitude,
			Title:     stops[i].BusStop.Description,
			Address:   fmt.Sprintf("%.0f m away", stops[i].Distance*1000),
			ReplyMarkup: ted.InlineKeyboardMarkup{
				InlineKeyboard: [][]ted.InlineKeyboardButton{
					{
						{
							Text: "Get ETAs",
							CallbackData: CallbackData{
								Type:   callbackNearbyETA,
								StopID: stops[i].BusStop.ID,
							}.Encode(),
						},
					},
				},
			},
		}
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(venues[0]).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(venues[1]).Return(ted.Response{}, nil)
	b.BusStops.(*MockBusStopRepository).EXPECT().Nearby(lat, lon, float32(1), 5).Return(stops, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Location: &ted.Location{
				Latitude:  lat,
				Longitude: lon,
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleReply_noLocationFound(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	var stops []NearbyBusStop
	var lat, lon float32 = 1.307574, 103.863256
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   stringNoLocationsNearby,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)
	b.BusStops.(*MockBusStopRepository).EXPECT().Nearby(lat, lon, float32(1), 5).Return(stops, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Location: &ted.Location{
				Latitude:  lat,
				Longitude: lon,
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}
