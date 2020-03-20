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
					Text:         callbackRefresh,
					CallbackData: "{\"t\":\"refresh\",\"b\":\"96049\",\"s\":[\"5\",\"24\"]}",
				},
				{
					Text:         callbackResend,
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

func Test_manageFavouritesReplyMarkup(t *testing.T) {
	expected := ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text:         "Add a new favourite",
					CallbackData: CallbackData{Type: callbackAddFavourite}.Encode(),
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

func Test_addFavouriteSuggestNameMarkup(t *testing.T) {
	t.Run("with bus stop description", func(t *testing.T) {
		expected := ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text:         "Opp Tropicana Condo",
						CallbackData: `{"t":"sf","n":"Opp Tropicana Condo","q":"96049 5 24"}`,
					},
				},
				{
					{
						Text:         "96049 5 24",
						CallbackData: `{"t":"sf","n":"96049 5 24","q":"96049 5 24"}`,
					},
				},
				{
					{
						Text:         "Set a custom name",
						CallbackData: "{\"t\":\"sf\",\"q\":\"96049 5 24\"}",
					},
				},
			},
		}
		actual := addFavouriteSuggestNameMarkup(Query{Stop: "96049", Filter: []string{"5", "24"}}, "Opp Tropicana Condo")
		assert.Equal(t, actual, expected)
	})
	t.Run("without bus stop description", func(t *testing.T) {
		expected := ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text:         "96049 5 24",
						CallbackData: `{"t":"sf","n":"96049 5 24","q":"96049 5 24"}`,
					},
				},
				{
					{
						Text:         "Set a custom name",
						CallbackData: "{\"t\":\"sf\",\"q\":\"96049 5 24\"}",
					},
				},
			},
		}
		actual := addFavouriteSuggestNameMarkup(Query{Stop: "96049", Filter: []string{"5", "24"}}, "")
		assert.Equal(t, actual, expected)
	})
}
