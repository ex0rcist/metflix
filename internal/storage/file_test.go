package storage

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/ex0rcist/metflix/internal/metrics"
)

func checkNoError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

func removeFile(t *testing.T, path string) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to remove file: %v", err)
	}
}

func TestNewFileStorage(t *testing.T) {
	storePath := "test_store.json"
	defer removeFile(t, storePath)

	fs, err := NewFileStorage(storePath, 0, false)
	checkNoError(t, err, "failed to create new FileStorage")

	if fs == nil {
		t.Fatalf("expected FileStorage instance, got nil")
	}
	if fs.storePath != storePath {
		t.Errorf("expected storePath=%s, got %s", storePath, fs.storePath)
	}
	if fs.storeInterval != 0 {
		t.Errorf("expected storeInterval=0, got %d", fs.storeInterval)
	}
	if fs.restoreOnStart {
		t.Errorf("expected restoreOnStart=false, got true")
	}
}

func TestSyncPushAndDump(t *testing.T) {
	ctx := context.Background()
	storePath := "test_store.json"
	defer removeFile(t, storePath)

	fs, err := NewFileStorage(storePath, 0, false)
	checkNoError(t, err, "failed to create new FileStorage")

	record := Record{Name: "test", Value: metrics.Counter(42)}
	err = fs.Push(ctx, record.CalculateRecordID(), record)
	checkNoError(t, err, "failed to push record")

	data, err := os.ReadFile(storePath)
	checkNoError(t, err, "failed to read storage file")

	err = json.Unmarshal(data, &fs.MemStorage)
	checkNoError(t, err, "failed to unmarshal storage file")

	if got, err := fs.Get(ctx, record.CalculateRecordID()); err != nil || got != record {
		t.Errorf("expected record %v, got %v", record, got)
	}
}

func TestRestore(t *testing.T) {
	ctx := context.Background()

	storePath := "test_store.json"
	defer removeFile(t, storePath)

	// create dump
	fs1, err := NewFileStorage(storePath, 0, false)
	checkNoError(t, err, "failed to create new FileStorage")

	record := Record{Name: "test", Value: metrics.Counter(42)}
	err = fs1.Push(ctx, record.CalculateRecordID(), record) // dumped
	checkNoError(t, err, "failed to push to FileStorage")

	// new storage from dump
	fs2, err := NewFileStorage(storePath, 0, true)
	checkNoError(t, err, "failed to create new FileStorage")

	restoredRecord, err := fs2.Get(ctx, record.CalculateRecordID())
	if err != nil {
		checkNoError(t, err, "expected to find restored record, but did not")
	}
	if restoredRecord != record {
		t.Errorf("expected restored record %v, got %v", record, restoredRecord)
	}
}

func TestAsyncDumping(t *testing.T) {
	ctx := context.Background()

	storePath := "test_store.json"
	defer removeFile(t, storePath)

	fs, err := NewFileStorage(storePath, 1, false)
	checkNoError(t, err, "failed to create new FileStorage")

	record := Record{Name: "test", Value: metrics.Counter(42)}
	err = fs.Push(ctx, record.CalculateRecordID(), record)
	checkNoError(t, err, "failed to push record")

	time.Sleep(1000 * time.Millisecond)

	err = fs.Close(ctx)
	checkNoError(t, err, "failed to close fs")

	data, err := os.ReadFile(storePath)
	checkNoError(t, err, "failed to read storage file")

	ms := NewMemStorage()
	err = json.Unmarshal(data, &ms)
	checkNoError(t, err, "failed to unmarshal storage file")

	if restored, err := ms.Get(ctx, record.CalculateRecordID()); err != nil || restored != record {
		t.Errorf("expected record %v, got %v", record, restored)
	}
}
