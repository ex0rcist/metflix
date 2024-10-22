package compression

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/klauspost/compress/gzip"
)

func TestCompressor_Write_SupportedContent(t *testing.T) {
	ctx := context.Background()
	recorder := httptest.NewRecorder()

	compressor := NewCompressor(recorder, ctx)
	compressor.Header().Set("Content-Type", "application/json")

	data := []byte(`{"message": "test"}`)
	n, err := compressor.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n != len(data) {
		t.Fatalf("expected %d bytes written, got %d", len(data), n)
	}

	compressor.Close()

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected Content-Encoding to be gzip, got %s", resp.Header.Get("Content-Encoding"))
	}

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("expected no error creating gzip reader, got %v", err)
	}
	defer gr.Close()

	uncompressedData := new(bytes.Buffer)
	_, err = uncompressedData.ReadFrom(gr)
	if err != nil {
		t.Fatalf("expected no error decompressing dara, got %v", err)
	}
	if !bytes.Equal(data, uncompressedData.Bytes()) {
		t.Fatalf("expected %s, got %s", data, uncompressedData.Bytes())
	}
}

func TestCompressor_Write_UnsupportedContent(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.Background()
	compressor := NewCompressor(recorder, ctx)
	compressor.Header().Set("Content-Type", "text/plain")

	data := []byte("test data")
	n, err := compressor.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if n != len(data) {
		t.Fatalf("expected %d bytes written, got %d", len(data), n)
	}

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		t.Fatalf("expected Content-Encoding to not be gzip")
	}

	body := recorder.Body.Bytes()
	if !bytes.Equal(data, body) {
		t.Fatalf("expected %s, got %s", data, body)
	}
}

func TestCompressor_Close(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context.Background()
	compressor := NewCompressor(recorder, ctx)
	compressor.Header().Set("Content-Type", "application/json")

	data := []byte(`{"message": "test"}`)
	_, err := compressor.Write(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	compressor.Close()

	if compressor.encoder != nil {
		t.Fatalf("expected encoder to be nil after close")
	}
}
