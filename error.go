package bete

import (
	"context"

	"github.com/getsentry/sentry-go"
)

func captureError(ctx context.Context, err error) {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	}
}
