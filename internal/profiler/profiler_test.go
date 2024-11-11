package profiler

import (
	"testing"
)

func TestGetProfilerSingleton(t *testing.T) {
	profiler1 := GetProfiler()
	profiler2 := GetProfiler()

	if profiler1 != profiler2 {
		t.Errorf("GetProfiler() должен возвращать один и тот же экземпляр, но получил разные")
	}
}
