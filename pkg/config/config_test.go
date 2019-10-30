package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	parameterPath = "/event-forwarder-tests/"
	sentryDSNParameter = "/event-forwarder-tests/sentry_dsn"

	cfg, err := New("foo")
	require.NoError(t, err)

	cases := []struct {
		got, want interface{}
	}{
		{cfg.Environment, "test"},
		{cfg.SentryDSN, "sentry_dsn"},
		{cfg.SentryTags, map[string]string{"function": "foo", "environment": "test"}},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.got)
	}
}
