package config

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// A Config provides a central location to store configs or credentials to
// instantiate other services.
type Config struct {
	Environment string
	SentryDSN   string
	SentryTags  map[string]string
}

const environmentENV = "ENVIRONMENT"

var (
	parameterPath      = "/event-forwarder/"
	sentryDSNParameter = "/event-forwarder/sentry_dsn"
)

// New returns a new instance of the config module corresponding to its specific
// environment.
func New(functionName string) (*Config, error) {
	// Authenticate with AWS
	sess, err := session.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "aws session error")
	}

	svc := ssm.New(sess)
	env := os.Getenv(environmentENV)

	cfg := &Config{
		Environment: env,
		SentryTags:  map[string]string{"function": functionName, "environment": env},
	}

	err = fetchParametersByPath(svc, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func fetchParametersByPath(svc *ssm.SSM, cfg *Config) error {
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(parameterPath),
		WithDecryption: aws.Bool(true),
	}

	resp, err := svc.GetParametersByPath(input)
	if err != nil {
		return errors.Wrap(err, "get parameters by path error")
	}

	for _, v := range resp.Parameters {
		if *v.Name == sentryDSNParameter {
			cfg.SentryDSN = *v.Value
			break
		}
	}

	return nil
}
