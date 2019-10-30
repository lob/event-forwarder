package handler

import (
	"context"

	"github.com/lob/event-forwarder/pkg/application"
)

// Handler contains the necessary fields and methods to successfully process a
// Lambda CloudWatch Logs event.
type Handler struct {
	app *application.App
}

// NewHandler returns a new Handler given an App.
func NewHandler(app *application.App) Handler {
	return Handler{app}
}

// Invoke is the main Lambda function handler.
func (h Handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	err := h.app.SQS.SendMessage(string(payload), h.app.Config.Tags)
	if err != nil {
		h.app.Sentry.CaptureAndWait(err)
	}

	return nil, err
}
