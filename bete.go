package bete

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/yi-jiayu/datamall/v3"
	"github.com/yi-jiayu/ted"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (c RealClock) Now() time.Time {
	return time.Now()
}

type DataMall interface {
	GetBusArrival(busStopCode string, serviceNo string) (datamall.BusArrival, error)
}

type Telegram interface {
	Do(request ted.Request) (ted.Response, error)
	DoMulti(requests ...ted.Request) ([]ted.Response, error)
}

type Bete struct {
	Version                string
	StreetViewStaticAPIKey string

	Clock      Clock
	BusStops   BusStopRepository
	Favourites FavouriteRepository
	DataMall   DataMall
	Telegram   Telegram
}

func sentrySetUser(ctx context.Context, id int) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{
				ID: strconv.Itoa(id),
			})
		})
	}
}

func (b Bete) HandleUpdate(ctx context.Context, u ted.Update) {
	switch {
	case u.Message != nil:
		sentrySetUser(ctx, u.Message.From.ID)
		telegramUpdatesTotal.WithLabelValues("message").Inc()

		b.HandleMessage(ctx, u.Message)
	case u.CallbackQuery != nil:
		sentrySetUser(ctx, u.CallbackQuery.From.ID)
		telegramUpdatesTotal.WithLabelValues("callback_query").Inc()

		b.HandleCallbackQuery(ctx, u.CallbackQuery)
	case u.InlineQuery != nil:
		telegramUpdatesTotal.WithLabelValues("inline_query").Inc()
		sentrySetUser(ctx, u.InlineQuery.From.ID)

		b.HandleInlineQuery(ctx, u.InlineQuery)
	case u.ChosenInlineResult != nil:
		telegramUpdatesTotal.WithLabelValues("chosen_inline_result").Inc()
		sentrySetUser(ctx, u.ChosenInlineResult.From.ID)

		b.HandleChosenInlineResult(ctx, u.ChosenInlineResult)
	}
}

// send makes a request to the Telegram Bot API.
func (b Bete) send(ctx context.Context, req ted.Request) {
	_, err := b.Telegram.Do(req)
	if err != nil {
		if ted.IsMessageNotModified(err) {
			return
		}
		captureError(ctx, errors.WithStack(err))
	}
}

func (b Bete) etaMessageText(ctx context.Context, stopID string, filter []string, format Format) (string, error) {
	t := b.Clock.Now()
	var errMsg string
	arrivals, err := b.DataMall.GetBusArrival(stopID, "")
	if err != nil {
		captureError(ctx, err)
		errMsg = "Error getting bus arrivals from LTA DataMall"
		if datamallErr, ok := err.(*datamall.Error); ok {
			errMsg += fmt.Sprintf(" (HTTP status %d)", datamallErr.StatusCode)
		}
	}
	var stop BusStop
	stop, err = b.BusStops.Find(stopID)
	if err != nil {
		if err != ErrNotFound {
			captureError(ctx, err)
		}
		stop = BusStop{ID: stopID}
	}
	info := ArrivalInfo{
		Stop:     stop,
		Time:     t,
		Services: arrivals.Services,
		Filter:   filter,
		ErrMsg:   errMsg,
	}
	return FormatArrivals(info, format)
}
