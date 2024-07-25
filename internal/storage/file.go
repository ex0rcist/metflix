package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ex0rcist/metflix/internal/logging"
)

// check that FileStorage implements MetricsStorage
var _ MetricsStorage = (*FileStorage)(nil)

type FileStorage struct {
	*MemStorage

	storePath string
	syncMode  bool
}

func NewFileStorage(storePath string, storeInterval int) *FileStorage {
	return &FileStorage{
		MemStorage: NewMemStorage(),
		storePath:  storePath,
		syncMode:   storeInterval == 0,
	}
}

func (s *FileStorage) Push(id string, record Record) error {
	if err := s.MemStorage.Push(id, record); err != nil {
		return err
	}

	if s.syncMode {
		return s.Dump()
	}

	return nil
}

func (s *FileStorage) Dump() (err error) {
	logging.LogInfo("dumping storage to file " + s.storePath)

	file, err := os.OpenFile(s.storePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("error during FileStorage.Dump()/os.OpenFile(): %w", err)
	}

	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("error during FileStorage.Dump()/file.Close(): %w", closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	snapshot := s.Snapshot()

	if err := encoder.Encode(snapshot); err != nil {
		return fmt.Errorf("error during FileStorage.Dump()/encoder.Encode(): %w", err)
	}

	return nil
}

func (s *FileStorage) Restore() (err error) {
	logging.LogInfo("restoring storage from file " + s.storePath)

	file, err := os.Open(s.storePath)
	if err != nil {
		if os.IsNotExist(err) {
			logging.LogWarn("no storage dump found to restore")
			return nil
		}

		return fmt.Errorf("error during FileStorage.Restore()/os.Open(): %w", err)
	}

	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = fmt.Errorf("error during FileStorage.Restore()/file.Close(): %w", closeErr)
		}
	}()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(s.MemStorage); err != nil {
		return fmt.Errorf("error during FileStorage.Restore()/decoder.Decode(): %w", err)
	}

	logging.LogInfo("storage data was restored")

	return nil
}

func (s *FileStorage) Kind() string {
	return KindFile
}
