package bete

import (
	"context"
	"fmt"

	"github.com/yi-jiayu/ted"
)

const resultsPerQuery = 50

func inlineQueryResult(stop BusStop, thumbURL string) ted.InlineQueryResult {
	return ted.InlineQueryResultArticle{
		ID:    stop.ID,
		Title: fmt.Sprintf("%s (%s)", stop.Description, stop.ID),
		InputMessageContent: ted.InputTextMessageContent{
			Text: fmt.Sprintf(`<strong>%s (%s)</strong>
%s
<pre>
Fetching ETAs...
</pre>`, stop.Description, stop.ID, stop.RoadName),
			ParseMode: "HTML",
		},
		ReplyMarkup: inlineETAMessageReplyMarkupP(stop.ID, FormatSummary),
		Description: stop.RoadName,
		ThumbURL:    thumbURL,
	}
}

func nearbyInlineQueryResult(stop NearbyBusStop, thumbURL string) ted.InlineQueryResult {
	return ted.InlineQueryResultArticle{
		ID:    stop.ID,
		Title: fmt.Sprintf("%s (%s)", stop.Description, stop.ID),
		InputMessageContent: ted.InputTextMessageContent{
			Text: fmt.Sprintf(`<strong>%s (%s)</strong>
%s
<pre>
Fetching ETAs...
</pre>`, stop.Description, stop.ID, stop.RoadName),
			ParseMode: "HTML",
		},
		ReplyMarkup: inlineETAMessageReplyMarkupP(stop.ID, FormatSummary),
		Description: fmt.Sprintf("%.0f m away", stop.Distance*1000),
		ThumbURL:    thumbURL,
	}
}

func (b Bete) searchBusStopsResults(query string) ([]ted.InlineQueryResult, error) {
	stops, err := b.BusStops.Search(query, resultsPerQuery)
	if err != nil {
		return nil, err
	}
	results := make([]ted.InlineQueryResult, len(stops))
	for i, stop := range stops {
		results[i] = inlineQueryResult(stop, getStreetViewStaticURL(b.StreetViewStaticAPIKey, stop))
	}
	return results, nil
}

func (b Bete) nearbyBusStopsResults(lat, lon float32) ([]ted.InlineQueryResult, error) {
	stops, err := b.BusStops.Nearby(lat, lon, 1, resultsPerQuery)
	if err != nil {
		return nil, err
	}
	results := make([]ted.InlineQueryResult, len(stops))
	for i, stop := range stops {
		results[i] = nearbyInlineQueryResult(stop, getStreetViewStaticURL(b.StreetViewStaticAPIKey, stop.BusStop))
	}
	return results, nil
}

func (b Bete) HandleInlineQuery(ctx context.Context, q *ted.InlineQuery) {
	var results []ted.InlineQueryResult
	var err error
	if q.Location != nil && q.Query == "" {
		results, err = b.nearbyBusStopsResults(q.Location.Latitude, q.Location.Longitude)
	} else {
		results, err = b.searchBusStopsResults(q.Query)
	}
	if err != nil {
		captureError(ctx, err)
		return
	}
	answer := ted.AnswerInlineQueryRequest{
		InlineQueryID: q.ID,
		Results:       results,
	}
	b.send(ctx, answer)
}

func (b Bete) HandleChosenInlineResult(ctx context.Context, r *ted.ChosenInlineResult) {
	stopID := r.ID
	text, err := b.etaMessageText(ctx, stopID, nil, FormatSummary)
	if err != nil {
		captureError(ctx, err)
		return
	}
	editMessageText := ted.EditMessageTextRequest{
		Text:            text,
		ParseMode:       "HTML",
		InlineMessageID: r.InlineMessageID,
		ReplyMarkup:     inlineETAMessageReplyMarkupP(stopID, FormatSummary),
	}
	b.send(ctx, editMessageText)
}
