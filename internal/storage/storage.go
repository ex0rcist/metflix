package storage

import "context"

// Kinds of storage
const (
	KindMemory   = "memory"
	KindFile     = "file"
	KindDatabase = "database"
)

// Common interface for storages: mem, file, db
type MetricsStorage interface {
	Push(ctx context.Context, id string, record Record) error
	PushList(ctx context.Context, data map[string]Record) error
	Get(ctx context.Context, id string) (Record, error)
	List(ctx context.Context) ([]Record, error)
	Close(ctx context.Context) error
}
