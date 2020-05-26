package bete

import (
	"context"
	"fmt"
	"testing"

	"github.com/yi-jiayu/ted"
)

func TestBete_handleAboutCommand(t *testing.T) {
	variants := []string{"/about", "/version"}
	for _, variant := range variants {
		t.Run(variant, func(t *testing.T) {
			b, finish := newMockBete(t)
			defer finish()

			version := randomStringID()
			userID := randomID()
			chatID := randomInt64ID()
			req := ted.SendMessageRequest{
				ChatID:    chatID,
				Text:      fmt.Sprintf(stringAboutMessage, version, version),
				ParseMode: "HTML",
			}

			b.Version = version
			b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

			update := ted.Update{
				Message: &ted.Message{
					ID:   randomID(),
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
				{
					{
						Text: "About Bus Eta Bot",
						CallbackData: CallbackData{
							Type: callbackAbout,
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

func TestBete_handleFavouritesCommand_nonPrivateChat(t *testing.T) {
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

func TestBete_handleBusStopCodeCommand(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	command := fmt.Sprintf("/%s", stop.ID)
	query := Query{Stop: stop.ID}
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

func TestBete_handleInvalidCommand(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	command := "/invalid"
	userID := randomID()
	chatID := randomInt64ID()
	req := ted.SendMessageRequest{
		ChatID: chatID,
		Text:   stringInvalidCommand,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
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
