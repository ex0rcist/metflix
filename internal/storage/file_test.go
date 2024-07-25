package storage

import (
	"os"
	"testing"

	"github.com/ex0rcist/metflix/internal/metrics"
	"github.com/stretchr/testify/require"
)

func createStoreWithData(t *testing.T, storePath string, storeInterval int) *FileStorage {
	store := NewFileStorage(storePath, storeInterval)

	err := store.Push("test_gauge", Record{Name: "test", Value: metrics.Gauge(42.42)})
	if err != nil {
		t.Fatalf("expected no error on push, got %v", err)
	}

	err = store.Push("test2_counter", Record{Name: "test2", Value: metrics.Counter(1)})
	if err != nil {
		t.Fatalf("expected no error on push, got %v", err)
	}

	return store
}

func TestSyncDumpRestoreStorage(t *testing.T) {
	storePath := "/tmp/db.json"

	t.Cleanup(func() {
		err := os.Remove(storePath)
		if err != nil {
			t.Fatalf("expected no error on cleanup, got %v", err)
		}
	})

	store := createStoreWithData(t, storePath, 0)
	storedData := store.Snapshot()

	store = NewFileStorage(storePath, 0)
	err := store.Restore()
	if err != nil {
		t.Fatalf("expected no error on restore, got %v", err)
	}

	restoredData := store.Snapshot()
	require.Equal(t, storedData, restoredData)
}

func TestAsyncDumpRestoreStorage(t *testing.T) {
	storePath := "/tmp/db.json"

	t.Cleanup(func() {
		err := os.Remove(storePath)
		if err != nil {
			t.Fatalf("expected no error on cleanup, got %v", err)
		}
	})

	store := createStoreWithData(t, storePath, 300)
	storedData := store.Snapshot()

	err := store.Dump()
	if err != nil {
		t.Fatalf("expected no error on dump, got %v", err)
	}

	store = NewFileStorage(storePath, 300)
	err = store.Restore()
	if err != nil {
		t.Fatalf("expected no error on restore, got %v", err)
	}

	restoredData := store.Snapshot()
	require.Equal(t, storedData, restoredData)
}

func TestRestoreDoesntFailIfNoSourceFile(t *testing.T) {
	store := NewFileStorage("xxx", 0)

	err := store.Restore()
	if err != nil {
		t.Fatalf("expected no error on empty file, got %v", err)
	}
}
