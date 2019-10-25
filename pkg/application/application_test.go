package application

import (
	"testing"

	"github.com/lob/event-forwarder/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{}
	app, err := New(cfg)
	require.NoError(t, err)

	assert.NotNil(t, app.Config)
	assert.NotNil(t, app.Sentry)
}
