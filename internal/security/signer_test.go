package security

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const secret = "abc"

func TestCalculateSignature(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		err      error
		expected string
	}{
		{
			name:     "Sign successful",
			data:     []byte("test data"),
			err:      nil,
			expected: "81173bfb4fd8a1691cf591b6fd17e72087d7efa1861227c10e8a493ec8d8c4f0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			signer := NewSignerService(secret)

			hash, err := signer.CalculateSignature(tt.data)

			require.ErrorIs(err, tt.err)
			require.Equal(tt.expected, hash)
		})
	}
}

func TestVerifySignature(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		hash  string
		valid bool
		err   error
	}{
		{
			name:  "Verify signature",
			data:  []byte("test data"),
			hash:  "81173bfb4fd8a1691cf591b6fd17e72087d7efa1861227c10e8a493ec8d8c4f0",
			valid: true,
			err:   nil,
		},

		{
			name:  "Bad signature",
			data:  []byte("test data"),
			hash:  "81173bfb4fd8a1691cf591b6fd17e72087d7efa1861227c10e8a493ec8d8c4f1",
			valid: false,
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			signer := NewSignerService(secret)

			valid, err := signer.VerifySignature(tt.data, tt.hash)

			require.ErrorIs(err, tt.err)
			require.Equal(tt.valid, valid)
		})
	}
}
