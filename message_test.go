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
	chatID := randomInt64ID()
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.Favourites.(*MockFavouriteRepository).EXPECT().Find(gomock.Any(), gomock.Any()).Return("")
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: randomID()},
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
	replyText := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        replyText,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
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
		Text:        fmt.Sprintf(AddFavouriteSuggestName, query.Canonical()),
		ReplyMarkup: addFavouriteSuggestNameMarkup(query, stop.Description),
	}

	b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: AddFavouritePromptForQuery,
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
		Text:        fmt.Sprintf(AddFavouriteSuggestName, query.Canonical()),
		ReplyMarkup: addFavouriteSuggestNameMarkup(query, ""),
	}

	b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(BusStop{}, ErrNotFound)
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: AddFavouritePromptForQuery,
			},
			Text: query.Canonical(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleReply_AddFavourite_Finish(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	query := "96049 5 24"
	favourites := []string{"Home", "SUTD"}
	name := "SUTD"
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        "New favourite added!",
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
				Text: fmt.Sprintf(AddFavouritePromptForName, query),
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
		Text:        AddFavouritePromptForQuery,
		ReplyMarkup: ted.ForceReply{},
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(reportError).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(askAgain).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			ReplyToMessage: &ted.Message{
				Text: AddFavouritePromptForQuery,
			},
			Text: messageText,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCommand_About(t *testing.T) {
	variants := []string{"/about", "/version"}
	for _, variant := range variants {
		t.Run(variant, func(t *testing.T) {
			b, finish := newMockBete(t)
			defer finish()

			version := randomStringID()
			userID := randomID()
			chatID := randomInt64ID()
			req := ted.SendMessageRequest{
				ChatID: chatID,
				Text:   "Bus Eta Bot " + version,
			}

			b.Version = version
			b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

			update := ted.Update{
				Message: &ted.Message{
					From: &ted.User{ID: userID},
					Chat: ted.Chat{ID: chatID, Type: "private"},
					Text: variant,
					Entities: []ted.MessageEntity{
						{
							Type:   "bot_command",
							Offset: 0,
							Length: len(variant),
						},
					},
				},
			}
			b.HandleUpdate(context.Background(), update)
		})
	}
}

func TestBete_HandleCommand_Favourite(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        "What would you like to do?",
		ReplyMarkup: manageFavouritesReplyMarkup(),
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID, Type: "private"},
			Text: "/favourites",
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 11,
				},
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCommand_Favourite_NonPrivateChat(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	req := ted.SendMessageRequest{
		ChatID:           chatID,
		Text:             "Sorry, you can only manage your favourites in a private chat.",
		ReplyToMessageID: messageID,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			ID:   messageID,
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID, Type: "group"},
			Text: "/favourites",
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 11,
				},
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func Test_getFavouriteQuery(t *testing.T) {
	text := fmt.Sprintf(AddFavouritePromptForName, "96049 5 24")
	assert.Equal(t, "96049 5 24", getFavouriteQuery(text))
	assert.Equal(t, "", getFavouriteQuery("invalid"))
}
