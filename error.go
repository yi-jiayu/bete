package bete

import (
	"context"
	"log"

	"github.com/getsentry/sentry-go"
)

func captureError(ctx context.Context, err error) {
	log.Printf("captured error: %v", err)
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	}
}

func captureMessage(ctx context.Context, msg string) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureMessage(msg)
	}
}
