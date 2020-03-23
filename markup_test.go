package bete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/ted"
)

func Test_etaMessageReplyMarkup(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected ted.InlineKeyboardMarkup
	}{
		{
			name:   "summary",
			format: FormatSummary,
			expected: ted.InlineKeyboardMarkup{
				InlineKeyboard: [][]ted.InlineKeyboardButton{
					{
						{
							Text:         stringFormatSwitchDetails,
							CallbackData: `{"t":"refresh","b":"96049","s":["5","24"],"f":"f"}`},
					},
					{
						{
							Text:         "Refresh",
							CallbackData: `{"t":"refresh","b":"96049","s":["5","24"],"f":"s"}`,
						},
						{
							Text:         "Resend",
							CallbackData: `{"t":"resend","b":"96049","s":["5","24"],"f":"s"}`,
						},
					},
				},
			},
		},
		{
			name:   "details",
			format: FormatDetails,
			expected: ted.InlineKeyboardMarkup{
				InlineKeyboard: [][]ted.InlineKeyboardButton{
					{
						{
							Text:         stringFormatSwitchSummary,
							CallbackData: `{"t":"refresh","b":"96049","s":["5","24"],"f":"s"}`,
						},
					},
					{
						{
							Text:         "Refresh",
							CallbackData: `{"t":"refresh","b":"96049","s":["5","24"],"f":"f"}`,
						},
						{
							Text:         "Resend",
							CallbackData: `{"t":"resend","b":"96049","s":["5","24"],"f":"f"}`,
						},
					},
				},
			},
		},
		{
			name: "defaults to summary",
			expected: ted.InlineKeyboardMarkup{
				InlineKeyboard: [][]ted.InlineKeyboardButton{
					{
						{
							Text:         stringFormatSwitchDetails,
							CallbackData: `{"t":"refresh","b":"96049","s":["5","24"],"f":"f"}`},
					},
					{
						{
							Text:         "Refresh",
							CallbackData: `{"t":"refresh","b":"96049","s":["5","24"],"f":"s"}`,
						},
						{
							Text:         "Resend",
							CallbackData: `{"t":"resend","b":"96049","s":["5","24"],"f":"s"}`,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := etaMessageReplyMarkup("96049", []string{"5", "24"}, tt.format)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_etaMessageReplyMarkupP(t *testing.T) {
	var markup ted.InlineKeyboardMarkup
	var markupP *ted.InlineKeyboardMarkup
	markup = etaMessageReplyMarkup("96049", []string{"5", "24"}, FormatSummary)
	markupP = etaMessageReplyMarkupP("96049", []string{"5", "24"}, FormatSummary)
	assert.Equal(t, markup, *markupP)
	markup = etaMessageReplyMarkup("96049", []string{"5", "24"}, FormatDetails)
	markupP = etaMessageReplyMarkupP("96049", []string{"5", "24"}, FormatDetails)
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

func Test_showFavouritesReplyMarkup_NoFavourites(t *testing.T) {
	expected := ted.ReplyKeyboardRemove{}
	actual := showFavouritesReplyMarkup(nil)
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
							Type: callbackDeleteFavourite,
							Name: "Home",
						}.Encode(),
					},
				},
				{
					{
						Text: "Work",
						CallbackData: CallbackData{
							Type: callbackDeleteFavourite,
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

func Test_inlineETAMessageReplyMarkupP(t *testing.T) {
	t.Run("summary", func(t *testing.T) {
		expected := &ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text:         stringFormatSwitchDetails,
						CallbackData: `{"t":"refresh","b":"96049","f":"f"}`},
				},
				{
					{
						Text:         "Refresh",
						CallbackData: `{"t":"refresh","b":"96049","f":"s"}`,
					},
				},
			},
		}
		actual := inlineETAMessageReplyMarkupP("96049", FormatSummary)
		assert.Equal(t, expected, actual)
	})
	t.Run("details", func(t *testing.T) {
		expected := &ted.InlineKeyboardMarkup{
			InlineKeyboard: [][]ted.InlineKeyboardButton{
				{
					{
						Text:         stringFormatSwitchSummary,
						CallbackData: `{"t":"refresh","b":"96049","f":"s"}`,
					},
				},
				{
					{
						Text:         "Refresh",
						CallbackData: `{"t":"refresh","b":"96049","f":"f"}`,
					},
				},
			},
		}
		actual := inlineETAMessageReplyMarkupP("96049", FormatDetails)
		assert.Equal(t, expected, actual)
	})
}
