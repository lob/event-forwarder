package sentry

import (
	"errors"
	"testing"

	"github.com/getsentry/raven-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRavenClient struct {
	packet *raven.Packet
}

func (c *mockRavenClient) Capture(packet *raven.Packet, tags map[string]string) (string, chan error) {
	c.packet = packet
	return "", make(chan error, 1)
}

func TestNew(t *testing.T) {
	client, err := New(Options{
		DSN:  "",
		Tags: map[string]string{"foo": "bar"},
	})
	require.NoError(t, err)

	assert.NotNil(t, client.rc)
}

func TestCaptureAndWait(t *testing.T) {
	rc := &mockRavenClient{}
	client := Sentry{rc}
	err := errors.New("test error")

	client.CaptureAndWait(err)

	ex, ok := rc.packet.Interfaces[0].(*raven.Exception)
	require.True(t, ok)

	assert.Equal(t, err.Error(), rc.packet.Message)
	assert.Equal(t, "TestCaptureAndWait", ex.Stacktrace.Frames[len(ex.Stacktrace.Frames)-2].Function)
}
