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

func TestBete_handleStartCommand(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	firstName := "Bete"
	command := "/start"
	version := randomStringID()
	userID := randomID()
	chatID := randomInt64ID()
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   fmt.Sprintf(stringWelcomeMessage, firstName),
		ReplyMarkup: ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text: "Take the tour!",
						CallbackData: CallbackData{
							Type: callbackTour,
							Name: tourSectionStart,
						}.Encode(),
					},
				},
			},
		},
	}

	b.Version = version
	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID, FirstName: firstName},
			Chat: ted.Chat{ID: chatID, Type: "private"},
			Text: command,
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: len(command),
				},
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_handleTourCommand(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	section := tour[tourSectionStart]
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   section.Text,
		ReplyMarkup: ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text: section.Navigation[0].Text,
						CallbackData: CallbackData{
							Type: callbackTour,
							Name: section.Navigation[0].Target,
						}.Encode(),
					},
				},
			},
		},
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: "/tour",
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 5,
				},
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_handleFavouritesCommand(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        stringFavouritesChooseAction,
		ReplyMarkup: favouritesReplyMarkup(),
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
		Text:             stringFavouritesOnlyPrivateChat,
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

func TestBete_handleETACommand_withArgs(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	command := "/eta"
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
			Chat: ted.Chat{ID: chatID, Type: "private"},
			Text: command + " " + query.Canonical(),
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: len(command),
				},
			},
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_handleETACommand_withInvalidArgs(t *testing.T) {
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

			command := "/eta"
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
					Chat: ted.Chat{ID: chatID, Type: "private"},
					Text: command + " " + tt.args,
					Entities: []ted.MessageEntity{
						{
							Type:   "bot_command",
							Offset: 0,
							Length: len(command),
						},
					},
				},
			}
			b.HandleUpdate(context.Background(), update)
		})
	}
}

func TestBete_handleETACommand_withoutArgs(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	command := "/eta"
	userID := randomID()
	messageID := randomID()
	chatID := randomInt64ID()
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        stringETACommandPrompt,
		ParseMode:   "HTML",
		ReplyMarkup: ted.ForceReply{},
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			ID:   messageID,
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: command,
			Entities: []ted.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: len(command),
				},
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
