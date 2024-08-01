package storage

import "context"

const (
	KindMemory   = "memory"
	KindFile     = "file"
	KindDatabase = "database"
)

// common interface for storages: mem, file, etc
type MetricsStorage interface {
	Push(ctx context.Context, id string, record Record) error
	PushList(ctx context.Context, data map[string]Record) error
	Get(ctx context.Context, id string) (Record, error)
	List(ctx context.Context) ([]Record, error)
	Close(ctx context.Context) error
}
