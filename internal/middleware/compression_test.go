package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecompressRequest_Success(t *testing.T) {
	data := []byte("test data")
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		t.Fatalf("expected no error writing writer.Write(), got %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decompressedData, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("expected no error reading decompressed data, got %v", err)
		}
		if !bytes.Equal(data, decompressedData) {
			t.Fatalf("expected %s, got %s", data, decompressedData)
		}
	})

	handler := DecompressRequest(nextHandler)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
}

func TestDecompressRequest_NoEncoding(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("test data")))

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "test data" {
			t.Fatalf("expected 'test data', got %s", string(body))
		}
	})

	handler := DecompressRequest(nextHandler)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
}

func TestDecompressRequest_UnsupportedEncoding(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("test data")))
	req.Header.Set("Content-Encoding", "deflate")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected handler not to be called")
	})

	handler := DecompressRequest(nextHandler)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestDecompressRequest_InternalError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid gzip data")))
	req.Header.Set("Content-Encoding", "gzip")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected handler not to be called")
	})

	handler := DecompressRequest(nextHandler)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func TestCompressResponse_Success(t *testing.T) {
	data := []byte("test data")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(data)

		if err != nil {
			t.Fatalf("expected no error writing writer.Write(), got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	handler := CompressResponse(nextHandler)

	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected Content-Encoding to be gzip, got %s", resp.Header.Get("Content-Encoding"))
	}

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("expected no error creating gzip reader, got %v", err)
	}
	defer gr.Close()

	decompressedData, err := io.ReadAll(gr)
	if err != nil {
		t.Fatalf("expected no error reading decompressed data, got %v", err)
	}

	if !bytes.Equal(data, decompressedData) {
		t.Fatalf("expected %s, got %s", data, decompressedData)
	}
}

func TestCompressResponse_NoCompressionRequested(t *testing.T) {
	data := []byte("test data")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(data)
		if err != nil {
			t.Fatalf("expected no error writing writer.Write(), got %v", err)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()
	handler := CompressResponse(nextHandler)

	handler.ServeHTTP(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("expected no error reading body, got %v", err)
	}

	if !bytes.Equal(data, body) {
		t.Fatalf("expected %s, got %s", data, body)
	}
}

func TestNeedGzipEncoding_Supported(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	if !needGzipEncoding(req) {
		t.Fatalf("expected needGzipEncoding to return true")
	}
}

func TestNeedGzipEncoding_Unsupported(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	if needGzipEncoding(req) {
		t.Fatalf("expected needGzipEncoding to return false")
	}
}
