package application

import (
	"github.com/lob/event-forwarder/pkg/config"
	"github.com/lob/event-forwarder/pkg/sentry"
	"github.com/lob/event-forwarder/pkg/sqs"
)

// App provides a central location to store references to clients
// or configurations to be used in the application.
type App struct {
	Config *config.Config
	Sentry *sentry.Sentry
	SQS    *sqs.SQS
}

// New returns a new instance of the application module.
func New(cfg *config.Config) (*App, error) {
	sentry, err := sentry.New(sentry.Options{
		DSN:  cfg.SentryDSN,
		Tags: cfg.Tags,
	})
	if err != nil {
		return nil, err
	}

	sqs, err := sqs.New(&sqs.Options{
		QueueURL: cfg.QueueURL,
	})
	if err != nil {
		return nil, err
	}

	return &App{
		Config: cfg,
		Sentry: sentry,
		SQS:    sqs,
	}, nil
}
