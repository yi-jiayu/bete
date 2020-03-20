package bete

import (
	"context"
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
}

type Bete struct {
	Version string

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
		b.HandleMessage(ctx, u.Message)
	case u.CallbackQuery != nil:
		b.HandleCallbackQuery(ctx, u.CallbackQuery)
	}
}

func (b Bete) etaMessageText(ctx context.Context, stopID string, filter []string) (string, error) {
	t := b.Clock.Now()
	arrivals, err := b.DataMall.GetBusArrival(stopID, "")
	if err != nil {
		return "", errors.Wrap(err, "error getting bus arrivals")
	}
	var stop BusStop
	stop, err = b.BusStops.Find(stopID)
	if err != nil {
		if err != ErrNotFound {
			captureError(ctx, err)
		}
		stop = BusStop{ID: stopID}
	}
	return FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     t,
		Services: arrivals.Services,
		Filter:   filter,
	})
}
