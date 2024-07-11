package logging_test

import (
	"testing"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	require := require.New(t)
	require.NotPanics(logging.Setup)
}
