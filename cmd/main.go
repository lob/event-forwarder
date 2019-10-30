package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lob/event-forwarder/pkg/application"
	"github.com/lob/event-forwarder/pkg/config"
	"github.com/lob/event-forwarder/pkg/handler"
	"github.com/lob/logger-go"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Err(err).Fatal("config error")
	}

	app, err := application.New(cfg)
	if err != nil {
		log.Err(err).Fatal("application error")
	}

	lambda.StartHandler(handler.NewHandler(app))
}
