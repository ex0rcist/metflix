package httpserver

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ex0rcist/metflix/internal/entities"
)

func TestNew(t *testing.T) {
	handler := http.NewServeMux()
	address := entities.Address("127.0.0.1:8080")

	server := New(handler, address)

	if server.server.Addr != address.String() {
		t.Errorf("expected server address %s, got %s", address.String(), server.server.Addr)
	}

	if server.server.Handler != handler {
		t.Errorf("expected server handler to be %v, got %v", handler, server.server.Handler)
	}

	if server.notify == nil {
		t.Error("expected notify channel to be initialized")
	}
}

func TestServerStart(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	address := entities.Address("127.0.0.1:8080")
	server := New(handler, address)

	server.Start()

	time.Sleep(100 * time.Millisecond)

	actualAddr := server.server.Addr

	req, err := http.NewRequest(http.MethodGet, "http://"+actualAddr+"/test", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request to server: %v", err)
	}
	defer func() {
		bcErr := resp.Body.Close()
		if bcErr != nil {
			t.Error(bcErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	err = server.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	select {
	case err, ok := <-server.Notify():
		if ok {
			if !errors.Is(err, http.ErrServerClosed) {
				t.Errorf("expected error 'http: Server closed', got %v", err)
			}
		}
	case <-time.After(time.Second):
		t.Error("timeout: notify channel did not close after shutdown")
	}
}
