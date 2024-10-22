package compression

import (
	"bytes"
	"fmt"

	"github.com/klauspost/compress/gzip"
)

// Pack []byte with gzip.
func Pack(data []byte) (*bytes.Buffer, error) {
	bb := new(bytes.Buffer)

	encoder, err := gzip.NewWriterLevel(bb, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	if _, err = encoder.Write(data); err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	if err = encoder.Close(); err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return bb, nil
}
