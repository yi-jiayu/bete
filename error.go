package bete

import (
	"context"
	"log"

	"github.com/getsentry/sentry-go"
)

func captureError(ctx context.Context, err error) {
	log.Printf("%+v", err)
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	}
}
