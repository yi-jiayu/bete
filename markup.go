package bete

import (
	"github.com/yi-jiayu/ted"
)

func etaMessageReplyMarkup(stopID string, filter []string, format Format) ted.InlineKeyboardMarkup {
	var rows [][]ted.InlineKeyboardButton
	if format == FormatDetails {
		rows = append(rows, []ted.InlineKeyboardButton{
			{
				Text: stringFormatSwitchSummary,
				CallbackData: CallbackData{
					Type:   callbackRefresh,
					StopID: stopID,
					Filter: filter,
					Format: FormatSummary,
				}.Encode(),
			},
		})
	} else {
		format = FormatSummary
		rows = append(rows, []ted.InlineKeyboardButton{
			{
				Text: stringFormatSwitchDetails,
				CallbackData: CallbackData{
					Type:   callbackRefresh,
					StopID: stopID,
					Filter: filter,
					Format: FormatDetails,
				}.Encode(),
			},
		})
	}
	rows = append(rows, []ted.InlineKeyboardButton{
		{
			Text: stringCallbackRefresh,
			CallbackData: CallbackData{
				Type:   callbackRefresh,
				StopID: stopID,
				Filter: filter,
				Format: format,
			}.Encode(),
		},
		{
			Text: stringCallbackResend,
			CallbackData: CallbackData{
				Type:   callbackResend,
				StopID: stopID,
				Filter: filter,
				Format: format,
			}.Encode(),
		},
	})
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func etaMessageReplyMarkupP(stopID string, filter []string, format Format) *ted.InlineKeyboardMarkup {
	markup := etaMessageReplyMarkup(stopID, filter, format)
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

func favouritesReplyMarkupP() *ted.InlineKeyboardMarkup {
	markup := favouritesReplyMarkup()
	return &markup
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
	rows = append(rows, []ted.InlineKeyboardButton{
		{
			Text: "Back",
			CallbackData: CallbackData{
				Type: callbackFavourites,
			}.Encode(),
		},
	})
	return &ted.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func inlineETAMessageReplyMarkupP(stopID string, format Format) *ted.InlineKeyboardMarkup {
	var showOtherFormat string
	var otherFormat Format
	if format == FormatDetails {
		showOtherFormat = stringFormatSwitchSummary
		otherFormat = FormatSummary
	} else {
		showOtherFormat = stringFormatSwitchDetails
		otherFormat = FormatDetails
	}
	return &ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text: showOtherFormat,
					CallbackData: CallbackData{
						Type:   callbackRefresh,
						StopID: stopID,
						Format: otherFormat,
					}.Encode(),
				},
			},
			{
				{
					Text: "Refresh",
					CallbackData: CallbackData{
						Type:   callbackRefresh,
						StopID: stopID,
						Format: format,
					}.Encode(),
				},
			},
		},
	}
}

func tourReplyMarkup(section TourSectionData) ted.ReplyMarkup {
	if len(section.Navigation) == 0 {
		return nil
	}
	var rows [][]ted.InlineKeyboardButton
	for _, nav := range section.Navigation {
		button := ted.InlineKeyboardButton{
			Text: nav.Text,
			CallbackData: CallbackData{
				Type: callbackTour,
				Name: string(nav.Target),
			}.Encode(),
		}
		rows = append(rows, []ted.InlineKeyboardButton{button})
	}
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
