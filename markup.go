package bete

import (
	"github.com/yi-jiayu/ted"
)

func etaMessageReplyMarkup(stopID string, filter []string) ted.InlineKeyboardMarkup {
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text: callbackRefresh,
					CallbackData: CallbackData{
						Type:   callbackRefresh,
						StopID: stopID,
						Filter: filter,
					}.Encode(),
				},
				{
					Text: callbackResend,
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
					Text: stringFavouritesManage,
					CallbackData: CallbackData{
						Type: callbackManageFavourites,
					}.Encode(),
				},
			},
			{
				{
					Text: stringFavouritesShow,
					CallbackData: CallbackData{
						Type: "show_favourites",
					}.Encode(),
				},
			},
			{
				{
					Text: stringFavouritesHide,
					CallbackData: CallbackData{
						Type: "hide_favourites",
					}.Encode(),
				},
			},
		},
	}
}

func showFavouritesReplyMarkup(favourites []string) ted.ReplyKeyboardMarkup {
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

func manageFavouritesReplyMarkupP(favourites []string) *ted.InlineKeyboardMarkup {
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
