package profiler

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ex0rcist/metflix/internal/logging"
)

var (
	profilerInstance *Profiler
	once             sync.Once
)

// Profiler service obj to make snapshots
type Profiler struct {
	callCounter int32
	createFile  func(name string) (*os.File, error)
}

// Singleton
func GetProfiler() *Profiler {
	once.Do(func() {
		profilerInstance = &Profiler{
			createFile: os.Create,
		}
	})
	return profilerInstance
}

// Save snapshot. Count is a special case to take snapshot in similar conditions.
func (p *Profiler) SaveMemoryProfile() {
	count := atomic.AddInt32(&p.callCounter, 1)

	if count != 3 {
		return
	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("./memory_profile_%s.prof", timestamp)

	f, err := p.createFile(fileName)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			logging.LogError(err)
		}
	}()

	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}

	fmt.Printf("Memory profile saved as %s\n", fileName)
}
