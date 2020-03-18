package bete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/ted"
)

func Test_etaMessageReplyMarkup(t *testing.T) {
	expected := ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text:         "Refresh",
					CallbackData: "{\"t\":\"refresh\",\"b\":\"96049\",\"s\":[\"5\",\"24\"]}",
				},
				{
					Text:         "Resend",
					CallbackData: "{\"t\":\"resend\",\"b\":\"96049\",\"s\":[\"5\",\"24\"]}",
				},
			},
		},
	}
	actual := etaMessageReplyMarkup("96049", []string{"5", "24"})
	assert.Equal(t, expected, actual)
}

func Test_etaMessageReplyMarkupP(t *testing.T) {
	markup := etaMessageReplyMarkup("96049", []string{"5", "24"})
	markupP := etaMessageReplyMarkupP("96049", []string{"5", "24"})
	assert.Equal(t, markup, *markupP)
}

func Test_favouritesReplyMarkup(t *testing.T) {
	expected := ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text:         "Add a new favourite",
					CallbackData: "{\"t\":\"add_favourite\"}",
				},
			},
			{
				{
					Text:         "Manage existing favourites",
					CallbackData: "{\"t\":\"edit_favourite\"}",
				},
			},
			{
				{
					Text:         "Show favourites keyboard",
					CallbackData: "{\"t\":\"show_favourites\"}",
				},
			},
			{
				{
					Text:         "Hide favourites keyboard",
					CallbackData: "{\"t\":\"hide_favourites\"}",
				},
			},
		},
	}
	actual := manageFavouritesReplyMarkup()
	assert.Equal(t, expected, actual)
}

func Test_showFavouritesReplyMarkup(t *testing.T) {
	expected := ted.ReplyKeyboardMarkup{
		Keyboard: [][]interface{}{
			{"Home"},
			{"Work"},
			{"MRT"},
		},
		ResizeKeyboard: true,
	}
	actual := showFavouritesReplyMarkup([]string{"Home", "Work", "MRT"})
	assert.Equal(t, expected, actual)
}
