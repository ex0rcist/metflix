package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ex0rcist/metflix/internal/logging"
	"github.com/ex0rcist/metflix/internal/utils"
)

// check that FileStorage implements MetricsStorage
var _ MetricsStorage = (*FileStorage)(nil)

type FileStorage struct {
	*MemStorage
	sync.Mutex

	storePath      string
	storeInterval  int
	restoreOnStart bool
	dumpTicker     *time.Ticker
}

func NewFileStorage(storePath string, storeInterval int, restoreOnStart bool) (*FileStorage, error) {
	fs := &FileStorage{
		MemStorage:     NewMemStorage(),
		storePath:      storePath,
		storeInterval:  storeInterval,
		restoreOnStart: restoreOnStart,
	}

	if fs.restoreOnStart {
		if err := fs.restore(); err != nil {
			return nil, err
		}
	}

	if fs.storeInterval > 0 {
		fs.dumpTicker = time.NewTicker(utils.IntToDuration(fs.storeInterval))
		go fs.startStorageDumping(fs.dumpTicker)
	}

	return fs, nil
}

func (s *FileStorage) Push(ctx context.Context, id string, record Record) error {
	if err := s.MemStorage.Push(ctx, id, record); err != nil {
		return err
	}

	if s.storeInterval == 0 {
		return s.dump()
	}

	return nil
}

func (s *FileStorage) PushList(ctx context.Context, data map[string]Record) error {
	if err := s.MemStorage.PushList(ctx, data); err != nil {
		return err
	}

	if s.storeInterval == 0 {
		return s.dump()
	}

	return nil
}

func (s *FileStorage) Close(_ context.Context) error {
	if s.dumpTicker != nil {
		s.dumpTicker.Stop()
	}

	return s.dump()
}

func (s *FileStorage) dump() (err error) {
	logging.LogInfo("dumping storage to file " + s.storePath)

	s.Lock()
	defer s.Unlock()

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

func (s *FileStorage) restore() (err error) {
	logging.LogInfo("restoring storage from file " + s.storePath)

	s.Lock()
	defer s.Unlock()

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

func (s *FileStorage) startStorageDumping(ticker *time.Ticker) {
	defer ticker.Stop()

	for {
		_, ok := <-ticker.C
		if !ok {
			break
		}

		if err := s.dump(); err != nil {
			logging.LogError(fmt.Errorf("error during FileStorage Dump(): %s", err.Error()))
		}
	}
}
