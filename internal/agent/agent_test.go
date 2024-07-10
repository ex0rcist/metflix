package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	_, err := New()
	require.NoError(t, err)
}
