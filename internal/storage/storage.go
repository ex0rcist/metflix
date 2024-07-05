package storage

type Storage interface {
	Push(record Record) error
	Get(recordID RecordID) (Record, error)
	GetAll() ([]Record, error)
}
