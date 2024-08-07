package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/ex0rcist/metflix/internal/storage"
	"github.com/stretchr/testify/mock"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "should return nil if storage online"},
		{name: "should return error if storage offline", err: entities.ErrUnexpected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := storage.NewPGXPoolMock()
			pm.On("Ping", mock.Anything).Return(tt.err)

			store := &storage.DatabaseStorage{Pool: pm}
			pinger := NewPingerService(store)

			err := pinger.Ping(context.Background())
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error to be %v, got: %v", tt.err, err)
			}

		})
	}
}

func TestPingOnUnpingableStorage(t *testing.T) {
	store := storage.NewMemStorage()
	pinger := NewPingerService(store)

	err := pinger.Ping(context.Background())
	if !errors.Is(err, entities.ErrStorageUnpingable) {
		t.Fatalf("expected error to be %v, got %v", entities.ErrStorageUnpingable, err)
	}
}
