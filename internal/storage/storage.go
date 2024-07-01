package storage

type Storage interface {
	Push(record Record) error
	Get(key string) (Record, error)
	GetAll() ([]Record, error)
}
