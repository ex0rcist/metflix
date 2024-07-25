package storage

const (
	KindMemory = "memory"
	KindFile   = "file"
	KindMock   = "mock"
)

// common interface for storages: mem, file, etc
type MetricsStorage interface {
	Push(id string, record Record) error
	Get(id string) (Record, error)
	List() ([]Record, error)
	Kind() string
}
