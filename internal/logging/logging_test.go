package logging

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	require := require.New(t)

	require.NotPanics(func() { Setup() })
}

func TestBasic(t *testing.T) {
	require := require.New(t)

	Setup()

	require.NotPanics(func() {
		LogInfo("some message")
		LogInfoCtx(context.Background(), "some message")
		LogInfoF("some message %d", 42)

		LogWarn("some message")
		LogWarnCtx(context.Background(), "some message")
		LogWarnF("some message %d", 42)

		LogError(errors.New("some message"))
		LogErrorCtx(context.Background(), errors.New("some message"))
		LogErrorF("some message %d", errors.New("test"))

		LogDebug("some message")
		LogDebugCtx(context.Background(), "some message")
		LogDebugF("some message %d", 42)
	})
}
