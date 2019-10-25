package sentry

import (
	"github.com/getsentry/raven-go"

	"github.com/pkg/errors"
)

type ravenClient interface {
	Capture(*raven.Packet, map[string]string) (string, chan error)
}

// Options provides fields to customize the creation of a Sentry struct.
type Options struct {
	DSN  string
	Tags map[string]string
}

// Sentry provides functions to capture errors and report them back to Sentry.
type Sentry struct {
	rc ravenClient
}

// New creates a new Sentry client.
func New(opts Options) (*Sentry, error) {
	rc, err := raven.NewWithTags(opts.DSN, opts.Tags)
	if err != nil {
		return nil, errors.Wrap(err, "sentry client error")
	}
	return &Sentry{rc}, nil
}

// CaptureAndWait takes in an error (that can optionally be wrapped by
// pkg/errors) and blocks while it reports it to Sentry.
func (s *Sentry) CaptureAndWait(err error) {
	stacktrace := raven.NewException(err, raven.GetOrNewStacktrace(err, 0, 2, nil))
	packet := raven.NewPacket(err.Error(), stacktrace)

	// the raven client doesn't have an AndWait helper function for packets
	eventID, ch := s.rc.Capture(packet, map[string]string{})
	if eventID != "" {
		<-ch
	}
}
