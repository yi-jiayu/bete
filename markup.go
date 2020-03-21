package bete

import (
	"github.com/yi-jiayu/ted"
)

func etaMessageReplyMarkup(stopID string, filter []string) ted.InlineKeyboardMarkup {
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text: stringCallbackRefresh,
					CallbackData: CallbackData{
						Type:   callbackRefresh,
						StopID: stopID,
						Filter: filter,
					}.Encode(),
				},
				{
					Text: stringCallbackResend,
					CallbackData: CallbackData{
						Type:   callbackResend,
						StopID: stopID,
						Filter: filter,
					}.Encode(),
				},
			},
		},
	}
}

func etaMessageReplyMarkupP(stopID string, filter []string) *ted.InlineKeyboardMarkup {
	markup := etaMessageReplyMarkup(stopID, filter)
	return &markup
}

func favouritesReplyMarkup() ted.InlineKeyboardMarkup {
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text: stringFavouritesAddNew,
					CallbackData: CallbackData{
						Type: callbackAddFavourite,
					}.Encode(),
				},
			},
			{
				{
					Text: stringFavouritesDelete,
					CallbackData: CallbackData{
						Type: callbackDeleteFavourites,
					}.Encode(),
				},
			},
			{
				{
					Text: stringFavouritesShow,
					CallbackData: CallbackData{
						Type: callbackShowFavourites,
					}.Encode(),
				},
			},
			{
				{
					Text: stringFavouritesHide,
					CallbackData: CallbackData{
						Type: callbackHideFavourites,
					}.Encode(),
				},
			},
		},
	}
}

func showFavouritesReplyMarkup(favourites []string) ted.ReplyMarkup {
	if len(favourites) == 0 {
		return ted.ReplyKeyboardRemove{}
	}
	var keyboard [][]interface{}
	for _, f := range favourites {
		keyboard = append(keyboard, []interface{}{f})
	}
	return ted.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func addFavouriteSuggestNameMarkup(query Query, description string) ted.InlineKeyboardMarkup {
	var rows [][]ted.InlineKeyboardButton
	if description != "" {
		rows = append(rows, []ted.InlineKeyboardButton{
			{
				Text: description,
				CallbackData: CallbackData{
					Type:  callbackSaveFavourite,
					Query: &query,
					Name:  description,
				}.Encode(),
			},
		})
	}
	rows = append(rows,
		[]ted.InlineKeyboardButton{
			{
				Text: query.Canonical(),
				CallbackData: CallbackData{
					Type:  callbackSaveFavourite,
					Query: &query,
					Name:  query.Canonical(),
				}.Encode(),
			},
		},
		[]ted.InlineKeyboardButton{
			{
				Text: stringFavouritesSetCustomName,
				CallbackData: CallbackData{
					Type:  callbackSaveFavourite,
					Query: &query,
				}.Encode(),
			},
		},
	)
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func deleteFavouritesReplyMarkupP(favourites []string) *ted.InlineKeyboardMarkup {
	if len(favourites) == 0 {
		return &ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text:         stringFavouritesAddNew,
						CallbackData: CallbackData{Type: callbackAddFavourite}.Encode(),
					},
				},
			},
		}
	}
	var rows [][]ted.InlineKeyboardButton
	for _, f := range favourites {
		rows = append(rows, []ted.InlineKeyboardButton{
			{
				Text: f,
				CallbackData: CallbackData{
					Type: callbackDeleteFavourite,
					Name: f,
				}.Encode(),
			},
		})
	}
	return &ted.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
