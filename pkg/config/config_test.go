package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	err := os.Setenv(queueENV, "localhost")
	require.NoError(t, err, "unexpected error setting queueENV value")
	err = os.Setenv(typeENV, "cloudwatch")
	require.NoError(t, err, "unexpected error setting typeENV value")
	defer func() {
		err := os.Setenv(queueENV, "")
		require.Nil(t, err, "unexpected error restoring original queueENV value")
		err = os.Setenv(typeENV, "")
		require.Nil(t, err, "unexpected error restoring original typeENV value")
	}()

	parameterPath = "/event-forwarder-tests/"
	sentryDSNParameter = "/event-forwarder-tests/sentry_dsn"

	cfg, err := New()
	require.NoError(t, err)

	cases := []struct {
		got, want interface{}
	}{
		{cfg.Environment, "test"},
		{cfg.SentryDSN, "sentry_dsn"},
		{cfg.Tags, map[string]string{"type": "cloudwatch", "environment": "test"}},
		{cfg.QueueURL, "localhost"},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.got)
	}
}
