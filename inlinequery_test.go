package bete

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/ted"
)

func Test_inlineQueryResult(t *testing.T) {
	stop := BusStop{
		ID:          "01319",
		Description: "Lavender Stn Exit A/ICA",
		RoadName:    "Kallang Rd",
		Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
	}
	expected := ted.InlineQueryResultArticle{
		ID:    "01319",
		Title: "Lavender Stn Exit A/ICA (01319)",
		InputMessageContent: ted.InputTextMessageContent{
			Text: `<strong>Lavender Stn Exit A/ICA (01319)</strong>
Kallang Rd
<pre>
Fetching ETAs...
</pre>`,
			ParseMode: "HTML",
		},
		ReplyMarkup: inlineETAMessageReplyMarkupP("01319", FormatSummary),
		Description: "Kallang Rd",
	}
	actual := inlineQueryResult(stop)
	assert.Equal(t, expected, actual)
}

func Test_nearbyInlineQueryResult(t *testing.T) {
	stop := NearbyBusStop{
		BusStop: BusStop{
			ID:          "01339",
			Description: "Bef Crawford Bridge",
			RoadName:    "Crawford St",
			Location:    Location{Latitude: 1.307746, Longitude: 103.864263},
		},
		Distance: 0.11356564947243729,
	}
	expected := ted.InlineQueryResultArticle{
		ID:    "01339",
		Title: "Bef Crawford Bridge (01339)",
		InputMessageContent: ted.InputTextMessageContent{
			Text: `<strong>Bef Crawford Bridge (01339)</strong>
Crawford St
<pre>
Fetching ETAs...
</pre>`,
			ParseMode: "HTML",
		},
		ReplyMarkup: inlineETAMessageReplyMarkupP("01339", FormatSummary),
		Description: "114 m away",
	}
	actual := nearbyInlineQueryResult(stop)
	assert.Equal(t, expected, actual)
}

func TestBete_HandleInlineQuery_Nearby(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stops := []NearbyBusStop{
		{
			BusStop: BusStop{
				ID:          "01339",
				Description: "Bef Crawford Bridge",
				RoadName:    "Crawford St",
				Location:    Location{Latitude: 1.307746, Longitude: 103.864263},
			},
			Distance: 0.11356564947243729,
		},
		{
			BusStop: BusStop{
				ID:          "07371",
				Description: "Aft Kallang Rd",
				RoadName:    "Lavender St",
				Location:    Location{Latitude: 1.309508, Longitude: 103.863501},
			},
			Distance: 0.21676780485189698,
		},
	}
	var lat, lon float32 = 1.307574, 103.863256
	query := ""
	inlineQueryID := randomStringID()
	answerInlineQuery := ted.AnswerInlineQueryRequest{
		InlineQueryID: inlineQueryID,
		Results: []ted.InlineQueryResult{
			nearbyInlineQueryResult(stops[0]),
			nearbyInlineQueryResult(stops[1]),
		},
	}

	b.BusStops.(*MockBusStopRepository).EXPECT().Nearby(lat, lon, float32(1), resultsPerQuery).Return(stops, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerInlineQuery).Return(ted.Response{}, nil)

	update := ted.Update{
		InlineQuery: &ted.InlineQuery{
			ID:       inlineQueryID,
			Location: &ted.Location{Latitude: lat, Longitude: lon},
			Query:    query,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleInlineQuery_Search(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stops := []BusStop{
		{
			ID:          "01319",
			Description: "Lavender Stn Exit A/ICA",
			RoadName:    "Kallang Rd",
			Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
		},
		{
			ID:          "07371",
			Description: "Aft Kallang Rd",
			RoadName:    "Lavender St",
			Location:    Location{Latitude: 1.309508, Longitude: 103.8635},
		},
	}
	query := "tropicana"
	inlineQueryID := randomStringID()
	answerInlineQuery := ted.AnswerInlineQueryRequest{
		InlineQueryID: inlineQueryID,
		Results: []ted.InlineQueryResult{
			inlineQueryResult(stops[0]),
			inlineQueryResult(stops[1]),
		},
	}

	b.BusStops.(*MockBusStopRepository).EXPECT().Search(query, resultsPerQuery).Return(stops, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerInlineQuery).Return(ted.Response{}, nil)

	update := ted.Update{
		InlineQuery: &ted.InlineQuery{
			ID:    inlineQueryID,
			Query: query,
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleInlineQuery_SearchNoResults(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	query := "tropicana"
	inlineQueryID := randomStringID()

	b.BusStops.(*MockBusStopRepository).EXPECT().Search(query, resultsPerQuery).Return(nil, nil)

	update := ted.Update{
		InlineQuery: &ted.InlineQuery{
			ID:    inlineQueryID,
			Query: query,
		},
	}
	b.HandleUpdate(context.Background(), update)
}
