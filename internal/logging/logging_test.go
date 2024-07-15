package logging_test

import (
	"context"
	"testing"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	require.NotPanics(func() { logging.Setup(ctx) })
}
