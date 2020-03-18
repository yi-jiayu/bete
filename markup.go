package bete

import (
	"github.com/yi-jiayu/ted"
)

func etaMessageReplyMarkup(stopID string, filter []string) ted.InlineKeyboardMarkup {
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text: "Refresh",
					CallbackData: CallbackData{
						Type:   "refresh",
						StopID: stopID,
						Filter: filter,
					}.Encode(),
				},
				{
					Text: "Resend",
					CallbackData: CallbackData{
						Type:   "resend",
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
					Text: "Add a new favourite",
					CallbackData: CallbackData{
						Type: "add_favourite",
					}.Encode(),
				},
			},
			{
				{
					Text: "Manage existing favourites",
					CallbackData: CallbackData{
						Type: "edit_favourite",
					}.Encode(),
				},
			},
			{
				{
					Text: "Show favourites keyboard",
					CallbackData: CallbackData{
						Type: "show_favourites",
					}.Encode(),
				},
			},
			{
				{
					Text: "Hide favourites keyboard",
					CallbackData: CallbackData{
						Type: "hide_favourites",
					}.Encode(),
				},
			},
		},
	}
}
