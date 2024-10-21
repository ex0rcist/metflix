package profiler

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"time"
)

type Profiler struct {
	callCounter int32
}

var (
	profilerInstance *Profiler
	once             sync.Once
)

func GetProfiler() *Profiler {
	once.Do(func() {
		profilerInstance = &Profiler{}
	})
	return profilerInstance
}

func (p *Profiler) SaveMemoryProfile() {
	count := atomic.AddInt32(&p.callCounter, 1)

	if count != 3 {
		return
	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("./memory_profile_%s.prof", timestamp)

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}

	fmt.Printf("Memory profile saved as %s\n", fileName)
}
