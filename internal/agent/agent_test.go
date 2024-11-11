package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	agnt, err := New()
	require.NoError(t, err)

	time.AfterFunc(5*time.Second, func() { agnt.Shutdown() })

	err = agnt.Run()
	require.NoError(t, err)
}
