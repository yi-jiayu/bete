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
					Text:         stringCallbackRefresh,
					CallbackData: "{\"t\":\"refresh\",\"b\":\"96049\",\"s\":[\"5\",\"24\"]}",
				},
				{
					Text:         stringCallbackResend,
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
					Text:         stringFavouritesAddNew,
					CallbackData: CallbackData{Type: callbackAddFavourite}.Encode(),
				},
			},
			{
				{
					Text:         stringFavouritesDelete,
					CallbackData: "{\"t\":\"delete_favourites\"}",
				},
			},
			{
				{
					Text:         stringFavouritesShow,
					CallbackData: "{\"t\":\"show_favourites\"}",
				},
			},
			{
				{
					Text:         stringFavouritesHide,
					CallbackData: "{\"t\":\"hide_favourites\"}",
				},
			},
		},
	}
	actual := favouritesReplyMarkup()
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
						Text:         stringFavouritesSetCustomName,
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
						Text:         stringFavouritesSetCustomName,
						CallbackData: "{\"t\":\"sf\",\"q\":\"96049 5 24\"}",
					},
				},
			},
		}
		actual := addFavouriteSuggestNameMarkup(Query{Stop: "96049", Filter: []string{"5", "24"}}, "")
		assert.Equal(t, actual, expected)
	})
}

func Test_deleteFavouritesReplyMarkup(t *testing.T) {
	t.Run("no favourites to delete", func(t *testing.T) {
		expected := &ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text:         stringFavouritesAddNew,
						CallbackData: CallbackData{Type: callbackAddFavourite}.Encode(),
					},
				},
			},
		}
		actual := deleteFavouritesReplyMarkupP(nil)
		assert.Equal(t, expected, actual)
	})
	t.Run("with favourites to delete", func(t *testing.T) {
		favourites := []string{"Home", "Work"}
		expected := &ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text: "Home",
						CallbackData: CallbackData{
							Type: callbackDeleteFavourites,
							Name: "Home",
						}.Encode(),
					},
				},
				{
					{
						Text: "Work",
						CallbackData: CallbackData{
							Type: callbackDeleteFavourites,
							Name: "Work",
						}.Encode(),
					},
				},
			},
		}
		actual := deleteFavouritesReplyMarkupP(favourites)
		assert.Equal(t, expected, actual)
	})
}
