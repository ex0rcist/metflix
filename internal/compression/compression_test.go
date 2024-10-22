package compression

import (
	"bytes"
	"io"
	"testing"

	"github.com/klauspost/compress/gzip"
)

func TestPack_Success(t *testing.T) {
	data := []byte("test data")
	buffer, err := Pack(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	reader, err := gzip.NewReader(buffer)
	if err != nil {
		t.Fatalf("expected no error creating gzip reader, got %v", err)
	}
	defer reader.Close()

	unpackedData, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("expected no error reading from gzip reader, got %v", err)
	}

	if !bytes.Equal(data, unpackedData) {
		t.Fatalf("expected %s, got %s", data, unpackedData)
	}
}
