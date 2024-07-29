package compression

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"

	"github.com/ex0rcist/metflix/internal/entities"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecompressor_Decompress_SupportedEncoding(t *testing.T) {
	data := []byte("test data")
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		t.Fatalf("expected no error on writer.Write(), got %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	ctx := context.Background()
	decompressor := NewDecompressor(req, ctx)

	err = decompressor.Decompress()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	decompressedData, err := io.ReadAll(decompressor.request.Body)
	if err != nil {
		t.Fatalf("expected no error reading decompressed data, got %v", err)
	}

	if !bytes.Equal(data, decompressedData) {
		t.Fatalf("expected %s, got %s", data, decompressedData)
	}

	decompressor.Close()
}

func TestDecompressor_Decompress_NoEncoding(t *testing.T) {
	data := []byte("test data")
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))

	ctx := context.Background()
	decompressor := NewDecompressor(req, ctx)

	err := decompressor.Decompress()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	decompressedData, err := io.ReadAll(decompressor.request.Body)
	if err != nil {
		t.Fatalf("expected no error reading data, got %v", err)
	}

	if !bytes.Equal(data, decompressedData) {
		t.Fatalf("expected %s, got %s", data, decompressedData)
	}

	decompressor.Close()
}

func TestDecompressor_Decompress_UnsupportedEncoding(t *testing.T) {
	data := []byte("test data")
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Set("Content-Encoding", "deflate")

	ctx := context.Background()
	decompressor := NewDecompressor(req, ctx)

	err := decompressor.Decompress()
	if !errors.Is(err, entities.ErrEncodingUnsupported) {
		t.Fatalf("expected %v, got %v", entities.ErrEncodingUnsupported, err)
	}
}

func TestDecompressor_Close(t *testing.T) {
	data := []byte("test data")
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		t.Fatalf("expected no error on writer.Write, got %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	ctx := context.Background()
	decompressor := NewDecompressor(req, ctx)

	err = decompressor.Decompress()
	if err != nil {
		t.Fatalf("expected no error on Decompress(), got %v", err)
	}

	decompressor.Close()

	if decompressor.reader != nil {
		t.Fatalf("expected reader to be nil after close")
	}
}

func TestDecompressor_Decompress_ErrorOnInit(t *testing.T) {
	// providing invalid gzip data to simulate error on NewReader
	data := []byte("invalid gzip data")
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	req.Header.Set("Content-Encoding", "gzip")

	ctx := context.Background()
	decompressor := NewDecompressor(req, ctx)

	err := decompressor.Decompress()
	if !errors.Is(err, entities.ErrEncodingInternal) {
		t.Fatalf("expected %v, got %v", entities.ErrEncodingInternal, err)
	}
}
